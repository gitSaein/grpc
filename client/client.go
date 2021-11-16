package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	rpc "grpc/rpc"

	"google.golang.org/grpc"
)

func main() {
	serverAddr := ":9988"
	cc, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fetchIt gRPC client failed tp dial to server: %v", err)
	}
	fc := rpc.NewFetchClient(cc)

	fIn := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, _, err := fIn.ReadLine()
		if err != nil {
			log.Printf("Failed to read a line in: %v", err)
		}

		ctx := context.Background()
		out, err := fc.Capitalize(ctx, &rpc.Payload{Data: line})
		if err != nil {
			log.Panicf("fetchIt gRPC Client got error from server: %v", err)
		}
		fmt.Printf("< %s\n\n", out.Data)
	}

}
