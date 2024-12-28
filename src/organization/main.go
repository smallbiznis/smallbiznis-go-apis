package main

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/smallbiznis/go-genproto/smallbiznis/balance/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-lib/pkg/env"
	"github.com/smallbiznis/go-lib/pkg/logger"
	"github.com/smallbiznis/go-lib/pkg/otelcol"
	"github.com/smallbiznis/go-lib/pkg/server"
	grpchandler "github.com/smallbiznis/organization/delivery/grpc"
	"github.com/smallbiznis/organization/infrastructure"
	"github.com/smallbiznis/organization/repository"
	"github.com/stripe/stripe-go/v80"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	stripe.Key = env.Lookup("STRIPE_SECRET_KEY", "sk_test_51LlDyWCIG3hXfWuZUvt4U1eu4ywFzTPj8eWePweMnXH9Bx96L1kCWbfmeFL0VsUI2TKgUUALxztYOcvbnHHyLyFB00fY5bYD4W")
}

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

func NewBalanceServiceClient() (balance.BalanceServiceClient, error) {
	conn, err := grpc.NewClient(env.Lookup("BALANCE_ADDR", ":4317"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return balance.NewBalanceServiceClient(conn), nil
}

func RegisterServiceServer(srv *grpc.Server, svc *grpchandler.OrganizationServiceSever) {
	organization.RegisterServiceServer(srv, svc)
}

func RegisterServiceHandlerFromEndpoint(mux *runtime.ServeMux) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return organization.RegisterServiceHandlerFromEndpoint(context.Background(), mux, env.Lookup("GRPC_PORT", ":4317"), opts)
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
	app := fx.New(
		fx.WithLogger(NewZapLogger),
		fx.Provide(infrastructure.NewGorm),
		fx.Invoke(
			Automigrate,
			Migrate,
		),
		otelcol.Resource,
		otelcol.TraceProvider,
		server.GrpcServerProvider,
		fx.Provide(NewBalanceServiceClient),
		fx.Provide(
			repository.NewCountryRepository,
			repository.NewOrganizationRepository,
			repository.NewLocationRepository,
			repository.NewRulesRepository,
			repository.NewShippingRateRepository,
			grpchandler.NewOrganizationServiceServer,
		),
		fx.Provide(NewServeMux, NewHttpServer),
		fx.Invoke(RegisterServiceServer, StartHTTPServer, RegisterServiceHandlerFromEndpoint),
		server.GrpcServerInvoke,
	)

	app.Run()
}
