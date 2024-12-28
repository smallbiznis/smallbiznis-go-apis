package main

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/smallbiznis/go-genproto/smallbiznis/member/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/organization/v1"
	"github.com/smallbiznis/go-genproto/smallbiznis/user/v1"
	"github.com/smallbiznis/go-lib/pkg/env"
	"github.com/smallbiznis/go-lib/pkg/logger"
	"github.com/smallbiznis/go-lib/pkg/otelcol"
	"github.com/smallbiznis/go-lib/pkg/server"
	grpchandler "github.com/smallbiznis/member/delivery/grpc"
	"github.com/smallbiznis/member/infrastructure"
	"github.com/smallbiznis/member/repository"
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

func RegisterServiceServer(srv *grpc.Server, svc *grpchandler.MemberServiceServer) {
	member.RegisterMemberServiceServer(srv, svc)
}

func RegisterServiceHandlerFromEndpoint(mux *runtime.ServeMux) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return member.RegisterMemberServiceHandlerFromEndpoint(context.Background(), mux, env.Lookup("GRPC_PORT", ":4317"), opts)
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
		return nil, err
	}

	return organization.NewServiceClient(conn), nil
}

func NewUserServiceClient() (user.ServiceClient, error) {
	conn, err := grpc.NewClient(env.Lookup("USER_ADDR", ":4317"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return user.NewServiceClient(conn), nil
}

func main() {
	app := fx.New(
		fx.WithLogger(NewZapLogger),
		fx.Provide(infrastructure.NewGorm),
		fx.Invoke(Automigrate),
		otelcol.Resource,
		otelcol.TraceProvider,
		server.GrpcServerProvider,
		fx.Provide(NewOrganizationServiceClient, NewUserServiceClient),
		fx.Provide(
			repository.NewOrganizationRepository,
			grpchandler.NewOrganizationServiceServer,
		),
		fx.Provide(NewServeMux, NewHttpServer),
		fx.Invoke(RegisterServiceServer, RegisterServiceHandlerFromEndpoint),
		fx.Invoke(StartHTTPServer),
		server.GrpcServerInvoke,
	)

	app.Run()
}
