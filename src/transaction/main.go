package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/smallbiznis/go-genproto/smallbiznis/customer/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/inventory/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/transaction/v1"
	"github.com/smallbiznis/go-lib/pkg/env"
	"github.com/smallbiznis/go-lib/pkg/logger"
	"github.com/smallbiznis/go-lib/pkg/otelcol"
	"github.com/smallbiznis/go-lib/pkg/server"
	grpchandler "github.com/smallbiznis/transaction/delivery/grpc"
	"github.com/smallbiznis/transaction/infrastructure"
	"github.com/smallbiznis/transaction/repository"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func NewOrganizationServiceClient() (organization.ServiceClient, error) {
	conn, err := grpc.NewClient(env.Lookup("ORGANIZATION_ADDR", ":4317"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return organization.NewServiceClient(conn), nil
}

func NewCustomerServiceClient() (customer.CustomerServiceClient, error) {
	conn, err := grpc.NewClient(env.Lookup("CUSTOMER_ADDR", ":4317"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return customer.NewCustomerServiceClient(conn), nil
}

func NewInventoryServiceClient() (inventory.ServiceClient, error) {
	conn, err := grpc.NewClient(env.Lookup("INVENTORY_ADDR", ":4317"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return inventory.NewServiceClient(conn), nil
}

func RegisterServiceServer(srv *grpc.Server, svc *grpchandler.TransactionService) {
	transaction.RegisterTransactionServiceServer(srv, svc)
}

func RegisterServiceHandlerFromEndpoint(mux *runtime.ServeMux) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return transaction.RegisterTransactionServiceHandlerFromEndpoint(context.Background(), mux, env.Lookup("GRPC_PORT", ":4317"), opts)
}

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := fx.New(
		fx.WithLogger(NewZapLogger),
		fx.Provide(infrastructure.NewGorm),
		fx.Invoke(Automigrate),
		otelcol.Resource,
		otelcol.TraceProvider,
		server.GrpcServerProvider,
		fx.Provide(
			NewOrganizationServiceClient,
			NewCustomerServiceClient,
			NewInventoryServiceClient,
		),
		fx.Provide(
			repository.NewOrderRepository,
			grpchandler.NewTransactionService,
		),
		fx.Provide(NewServeMux, NewHttpServer),
		fx.Invoke(
			RegisterServiceServer,
			StartHTTPServer,
			RegisterServiceHandlerFromEndpoint,
		),
		server.GrpcServerInvoke,
	)

	go func() {
		if err := app.Start(ctx); err != nil {
			zap.L().Fatal("app.Start", zap.Error(err))
		}
	}()

	<-ctx.Done()
	if err := app.Stop(context.TODO()); err != nil {
		zap.Error(err)
	}
}
