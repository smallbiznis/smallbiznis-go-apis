package main

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/smallbiznis/go-lib/pkg/env"
	"github.com/smallbiznis/go-lib/pkg/logger"
	"github.com/smallbiznis/go-lib/pkg/otelcol"
	"github.com/smallbiznis/go-lib/pkg/server"
	"github.com/smallbiznis/notification/infrastructure"
	"github.com/smallbiznis/notification/repository"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func NewZapLogger() fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger.InitLogger()}
}

func NewServeMux() *runtime.ServeMux {
	return runtime.NewServeMux()
}

func NewHttpServer(mux *runtime.ServeMux) *http.Server {
	return &http.Server{
		Addr:    env.Lookup("HTTP_PORT", ":4318"),
		Handler: mux,
	}
}

// func RegisterServiceServer(srv *grpc.Server, svc *grpchandler.OrganizationServiceSever) {
// 	organization.RegisterServiceServer(srv, svc)
// }

// func RegisterServiceHandlerFromEndpoint(mux *runtime.ServeMux) error {
// 	opts := []grpc.DialOption{
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	}
// 	return organization.RegisterServiceHandlerFromEndpoint(context.Background(), mux, env.Lookup("GRPC_PORT", ":4317"), opts)
// }

func StartHTTPServer(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}

func main() {
	app := fx.New(
		fx.WithLogger(NewZapLogger),
		fx.Provide(infrastructure.NewGorm),
		fx.Invoke(Automigrate, Migrate),
		otelcol.Resource,
		otelcol.TraceProvider,
		server.GrpcServerProvider,
		fx.Provide(
			repository.NewNotificaionRepository,
		),
		fx.Provide(NewServeMux, NewHttpServer),
		fx.Invoke(
			// RegisterServiceServer,
			StartHTTPServer,
			// RegisterServiceHandlerFromEndpoint,
		),
		server.GrpcServerInvoke,
	)

	app.Run()
}
