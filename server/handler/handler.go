package handler

import (
	"context"
	pb "grpc/build/proto/api"
	"log"
)

type APIServer struct {
}

func (s *APIServer) GetHello(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	log.Printf("Received: %v", in.GetName())

	return &pb.Reply{Message: "Hello " + in.GetName()}, nil
}
