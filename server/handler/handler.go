package handler

import (
	"context"
	"log"

	pb "grpc/build/proto/api"
)

// APIServer is representation of protobuf ApiServer
type APIServer struct {
	pb.UnimplementedApiServer
}

// GetHello implements api.proto.ApiServer.GetHello
func (s *APIServer) GetHello(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.Reply{Message: "Hello " + in.GetName()}, nil
}
