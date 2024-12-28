package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nfnt/resize"
	"github.com/smallbiznis/go-lib/pkg/env"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	serviceName  = env.Lookup("SERVICE_NAME", "storagesvc")
	endpoint     = env.Lookup("MINIO_ADDR", "localhost:9000")
	accessKey    = env.Lookup("MINIO_ACCESS_KEY", "HjzDQqc6xnStZgx87YV9")
	secretKey    = env.Lookup("MINIO_SECRET_KEY", "89WA1joQsK76nDnld8g7rYWzfD6EPt1RinwJyHQe")
	useSsl       = env.Lookup("MINIO_USE_SSL", "false")
	client       *minio.Client
	otelendpoint = env.Lookup("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318")
	tracer       = otel.Tracer(serviceName)
)

func initializeMinio() {
	val, err := strconv.ParseBool(useSsl)
	if err != nil {
		log.Fatal(err)
	}

	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: val,
	})
	if err != nil {
		log.Fatal(err)
	}

	client = c
}

func init() {
	initializeMinio()
}

func initResource(ctx context.Context) (*resource.Resource, error) {
	extraResources, err := resource.New(ctx,
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(env.Lookup("SERVICE_NAME", "storagesvc")),
			semconv.ServiceVersion(env.Lookup("SERVICE_VERSION", "1.0.0")),
			semconv.ServiceNamespace(env.Lookup("SERVICE_NAMESPACE", "example")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	resource, _ := resource.Merge(
		resource.Default(),
		extraResources,
	)

	return resource, nil
}

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initTraceProvider(ctx context.Context) (func(context.Context) error, error) {

	resource, _ := initResource(ctx)

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Set up a trace exporter
	// tracerExp, err := stdouttrace.New()
	// if err != nil {
	// 	zap.Error(err)
	// }

	// HTTP Exporter
	traceClient := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(otelendpoint),
	)

	tracerExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		zap.Error(err)
		return nil, err
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(tracerExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func initMetricProvider(ctx context.Context) (func(context.Context) error, error) {

	resource, _ := initResource(ctx)

	// Set up a metrics exporter
	metricClient, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(otelendpoint),
	)
	if err != nil {
		return nil, err
	}

	mp := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(
			metric.NewPeriodicReader(metricClient),
		),
	)
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetMeterProvider(mp)

	return mp.Shutdown, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize TraceProvider
	shutdown, err := initTraceProvider(ctx)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err.Error())
		}
	}()

	app := gin.New()

	v1 := app.Group("v1")

	buckets := v1.Group("buckets")
	{
		buckets.GET("", func(c *gin.Context) {
			buckets, err := client.ListBuckets(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			c.JSON(http.StatusOK, buckets)
		})
	}

	bucket := buckets.Group("/:bucket")
	{
		bucket.POST("", func(c *gin.Context) {
			ctx := c.Request.Context()

			bucket := c.Param("bucket")
			ok, err := client.BucketExists(ctx, bucket)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			if ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": gin.H{
						"message": "Bucket already exist",
					},
				})
				return
			}

			if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{ObjectLocking: true}); err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			c.Status(http.StatusOK)
		})

		bucket.PUT("/*object", func(c *gin.Context) {
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, c.Request.Body); err != nil {
				c.String(http.StatusInternalServerError, "Failed to read request body: %s", err.Error())
				return
			}

			fileBytes := buf.Bytes()
			fileSize := int64(len(fileBytes))

			if _, err := client.PutObject(c.Request.Context(), c.Param("bucket"), c.Param("object"), bytes.NewReader(fileBytes), fileSize, minio.PutObjectOptions{
				ContentType: "application/octet-stream",
			}); err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			c.Status(http.StatusOK)
		})

		bucket.GET("/*object", func(c *gin.Context) {
			ctx, span := tracer.Start(ctx, c.Request.URL.Path)
			defer span.End()

			TestSpan(ctx)

			bucket := c.Param("bucket")
			object := c.Param("object")
			obj, err := client.GetObject(ctx, bucket, object, minio.GetObjectOptions{})
			if err != nil {
				// Handle case where object is not found
				minioErr, ok := err.(minio.ErrorResponse)
				if ok && minioErr.StatusCode == http.StatusNotFound {
					c.JSON(http.StatusNotFound, gin.H{
						"error": "Object not found",
					})
					return
				}

				// Handle other potential errors
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to retrieve object",
				})
				return
			}
			defer obj.Close()

			// Mengambil nama file dari path
			fileName := filepath.Base(object)

			// Read the object content
			stat, err := obj.Stat()
			if err != nil {
				// Handle case where object stat fails (object might not exist)
				minioErr, ok := err.(minio.ErrorResponse)
				if ok && minioErr.StatusCode == http.StatusNotFound {
					c.JSON(http.StatusNotFound, gin.H{
						"error": "Object not found",
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Unable to stat object",
					})
				}
				return
			}

			// Serve the object as a file download
			c.Header("Content-Disposition", "attachment; filename="+fileName)
			c.Header("Content-Type", stat.ContentType)

			width, _ := strconv.Atoi(c.Query("width"))
			height, _ := strconv.Atoi(c.Query("height"))
			if width > 0 && height > 0 {
				img, _, err := image.Decode(obj)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
				}

				newImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

				// Encode gambar hasil resize ke buffer
				var buf bytes.Buffer
				if err := png.Encode(&buf, newImage); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image"})
					return
				}

				// Set header untuk mengirim gambar sebagai octet-stream
				c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
				return
			}

			c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))
			// Stream the file content to the response
			if _, err := io.Copy(c.Writer, obj); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error streaming the object",
				})
			}
		})
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}

// func ImageProcessing(obj *minio.Object, stat minio.ObjectInfo, width, height int) {
// 	dst := image.NewRGBA(image.Rect(0, 0, width, height))
// 	png.Encode(obj, dst)
// }

func TestSpan(ctx context.Context) {
	spanCtx := trace.SpanContextFromContext(ctx)
	fmt.Printf("TraceID: %v\n", spanCtx.TraceID())
	fmt.Printf("SpanID: %v\n", spanCtx.SpanID())

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetName("TestSpan")
}
