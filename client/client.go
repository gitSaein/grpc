package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "grpc/build/proto/api"
)

const (
	address = "127.0.0.1:50051"
	name    = "AaronRoh"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewApiClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	reply, err := c.GetHello(ctx, &pb.Request{Name: name})

	if err != nil {
		log.Fatalf("GetHello error: %v", err)
	}
	log.Printf("Person: %v", reply)
}
