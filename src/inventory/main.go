package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/smallbiznis/go-genproto/smallbiznis/inventory/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/env"
	"github.com/smallbiznis/go-lib/pkg/logger"
	"github.com/smallbiznis/go-lib/pkg/otelcol"
	"github.com/smallbiznis/go-lib/pkg/server"
	grpchandler "github.com/smallbiznis/inventory/delivery/grpc"
	"github.com/smallbiznis/inventory/infrastructure"
	"github.com/smallbiznis/inventory/repository"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
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

func RegisterServiceServer(srv *grpc.Server, svc *grpchandler.InventoryService) {
	inventory.RegisterServiceServer(srv, svc)
}

func RegisterServiceHandlerFromEndpoint(mux *runtime.ServeMux) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return inventory.RegisterServiceHandlerFromEndpoint(context.Background(), mux, env.Lookup("GRPC_PORT", ":4317"), opts)
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

func NewOrganizationServiceClient() (organization.ServiceClient, error) {
	conn, err := grpc.NewClient(env.Lookup("ORGANIZATION_ADDR", ":4317"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	return organization.NewServiceClient(conn), nil
}

func main() {
	app := fx.New(
		fx.WithLogger(NewZapLogger),
		fx.Provide(infrastructure.NewGorm, infrastructure.NewElastic),
		fx.Invoke(Automigrate),
		otelcol.Resource,
		otelcol.TraceProvider,
		server.GrpcServerProvider,
		fx.Provide(NewOrganizationServiceClient),
		fx.Provide(
			repository.NewInventoryRepository,
			grpchandler.NewInventoryService,
		),
		fx.Provide(NewServeMux, NewHttpServer),
		fx.Invoke(RegisterServiceServer, StartHTTPServer, RegisterServiceHandlerFromEndpoint),
		server.GrpcServerInvoke,
	)

	app.Run()
}
