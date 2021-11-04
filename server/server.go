package main

import (
	"log"
	"net"

	pb "grpc/build/api"
	handler "grpc/server/handler"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterApiServer(grpcServer, &handler.APIServer{})
	//추가

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
