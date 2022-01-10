package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "grpc/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	interceptor "grpc/interceptors"
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

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	// deadline
	for i := 0; i < 5; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "HelloworldService.SayHello canceled")
		}

		time.Sleep(1 * time.Second)
	}

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func ListenAndGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor))

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
