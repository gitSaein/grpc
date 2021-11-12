package middleware

import (
	"context"
	"grpc/server/handler"
	"net"

	"google.golang.org/grpc"
	"honnef.co/go/tools/config"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

type AppServer struct {
	app.AppServer
	cfg config.Config
}

func NewAppServer(cfg config.Config) (*AppServer, error) {
	return &AppServer{cfg: cfg}, nil
}

func (s *AppServer) HealthCheck(
	ctx context.Context,
	req *app.HealthCheckRequest,
) (*app.HealthCheckResponse, error) {
	return handler.HealthCheck()(ctx, req)
}

func NewGRPCServer(cfg config.Config) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
		),
	)

	appServer, err := NewAppServer(cfg)
	if err != nil {
		return nil, err
	}

	app.RegisterApiServer(grpcServer, appServer)

	return grpcServer, nil
}

func ServerGRPC(cfg config.Config) error {
	lis, err := net.Listen("tcp", ":"+cfg.Setting().GRPCServerPort)
	if err != nil {
		return err
	}

	grpcServer, err := NewGRPCServer(cfg)
	if err != nil {
		return err
	}

	return grpcServer.Serve(lis)
}
