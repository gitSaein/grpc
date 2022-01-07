package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	pb "gitlab.bemilycorp.com/prototype/echo-grpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

type Server struct {
	proto           string
	addr            string
	networkListener net.Listener
}

func getHttpHeader(ctx context.Context) {

	log.Println("================== HEADER start ==================")

	method, _ := grpc.Method(ctx)
	log.Printf("method: %s", method)
	md, _ := metadata.FromIncomingContext(ctx)

	log.Printf("metadata: %v", md)
	log.Println("==================  HEADER end  ==================")

}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	start := time.Now()

	getHttpHeader(ctx)

	// deadline
	for i := 0; i < 5; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "HelloworldService.SayHello canceled")
		}

		time.Sleep(1 * time.Second)
	}

	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))

	elapsed := time.Since(start)
	log.Printf("Received: %v, take time - %s", in.GetName(), elapsed)

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func ListenAndGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{}) // helloworld_grpc.pb.go 에 있음
	reflection.Register(s)                 // grpcurl 명령을 사용하게 하기 위해

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //grpc 서버 시작
		log.Fatalf("failed to serve: %v", err)
	}

}

func main() {

	ListenAndGrpcServer()

}
