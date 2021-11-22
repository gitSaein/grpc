package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	pb "grpc/prometheus/proto"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DemoServiceServer defines a Server.
type DemoServiceServer struct {
	pb.UnimplementedDemoServiceServer
}

func newDemoServer() *DemoServiceServer {
	return &DemoServiceServer{}
}

// SayHello implements a interface defined by protobuf.
func (s *DemoServiceServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	customizedCounterMetric.WithLabelValues(request.Name).Inc()
	// log.Fatalf("SayHello 함수: %v", request.Name)
	log.Printf("message send: %v", string(request.Name))

	return &pb.HelloResponse{Message: fmt.Sprintf("Hello %s", request.Name)}, nil
}

var (
	// Create a metrics registry.
	reg = prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Create a customized counter metric.
	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

func init() {
	// Register standard server metrics and customized metrics to registry.
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues("Test")
}

// NOTE: Graceful shutdown is missing. Don't use this demo in your production setup.
func main() {
	// Listen an actual port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9093))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Create a HTTP server for prometheus.
	// httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9092)}
	httpServer := &http.Server{Handler: promhttp.Handler(), Addr: fmt.Sprintf("0.0.0.0:%d", 9092)}

	// Create a gRPC Server with gRPC interceptor.
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	// Create a new api server.
	demoServer := newDemoServer()

	// Register your service.
	pb.RegisterDemoServiceServer(grpcServer, demoServer)

	// Initialize all metrics.
	grpcMetrics.InitializeMetrics(grpcServer)

	// Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Start your gRPC server.
	log.Fatal(grpcServer.Serve(lis))
}
