package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "grpc/prometheus/proto"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Create a metrics registry.
	reg := prometheus.NewRegistry()
	// Create some standard client metrics.
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	// Register client metrics to registry.
	reg.MustRegister(grpcMetrics)
	// Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%v", 9093),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// Create a HTTP server for prometheus.
	// httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9094)}
	httpServer := &http.Server{Handler: promhttp.Handler(), Addr: fmt.Sprintf("0.0.0.0:%d", 9094)}

	// Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Create a gRPC server client.
	client := pb.NewDemoServiceClient(conn)
	fmt.Println("Start to call the method called SayHello every 3 seconds")

	// ==========================================================================================

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		wg.Done()
		k := 0
		for i := 0; i < 3; i++ {
			// Call “SayHello” method and wait for response from gRPC Server.
			k++
			sliceMsg := []string{"[guest 1] ", strconv.Itoa(k)}
			message := strings.Join(sliceMsg, " - ")
			log.Printf("message send: %v", message)

			_, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: message})
			if err != nil {
				log.Printf("Calling the SayHello method unsuccessfully. ErrorInfo: %+v", err)
				log.Printf("You should to stop the process")
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()
	go func() {
		wg.Done()
		k := 0
		for i := 0; i < 3; i++ {
			// Call “SayHello” method and wait for response from gRPC Server.
			k++
			sliceMsg := []string{"[guest 2] ", strconv.Itoa(k)}
			message := strings.Join(sliceMsg, " - ")
			log.Printf("message send: %v", message)

			_, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: message})
			if err != nil {
				log.Printf("Calling the SayHello method unsuccessfully. ErrorInfo: %+v", err)
				log.Printf("You should to stop the process")
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		wg.Done()
		k := 0
		for i := 0; i < 3; i++ {
			// Call “SayHello” method and wait for response from gRPC Server.
			k++
			sliceMsg := []string{"[guest 3] ", strconv.Itoa(k)}
			message := strings.Join(sliceMsg, " - ")
			log.Printf("message send: %v", message)

			_, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: message})
			if err != nil {
				log.Printf("Calling the SayHello method unsuccessfully. ErrorInfo: %+v", err)
				log.Printf("You should to stop the process")
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()
	wg.Wait()

	// ==========================================================================================

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("You can press n or N to stop the process of client")
	for scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "n" {
			os.Exit(0)
		}
	}
}
