package handler

import (
	"context"
	"log"

	pb "grpc/build/api"
)

// APIServer is representation of protobuf ApiServer
type APIServer struct {
}

// GetHello implements api.proto.ApiServer.GetHello
func (s *APIServer) GetHello(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("Received: %v", in.GetName())

	return &pb.Response{Message: "Hello " + in.GetName()}, nil
}
