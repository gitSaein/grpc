package main

import (
	"bytes"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	rpc "grpc/rpc"
)

type fetchIt int

// Compile time assertion that fetchIt implements FetchServer.
var _ rpc.FetchServer = (*fetchIt)(nil)

func (fi *fetchIt) Capitalize(ctx context.Context, in *rpc.Payload) (*rpc.Payload, error) {
	out := &rpc.Payload{
		Data: bytes.ToUpper(in.Data),
	}
	return out, nil
}

func main() {
	addr := ":9988"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("gRPC server: failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	rpc.RegisterFetchServer(srv, new(fetchIt))
	log.Printf("fetchIt gRPC server serving at %q", addr)
	if err := srv.Serve(ln); err != nil {
		log.Fatalf("gRPC server: error serving: %v", err)
	}
}
