package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	localEndpointURI := "192.168.1.5:5454"
	reporterURI := "http://localhost:9411/api/v2/spans"
	serviceName := "server"
	// 1. zipkin exporter 설정
	localEndpoint, err := openzipkin.NewEndpoint(serviceName, localEndpointURI)
	if err != nil {
		log.Fatalf("Failed to create the local zipkinEndpoint: %v", err)
	}
	reporter := zipkinHTTP.NewReporter(reporterURI)
	ze := zipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(ze)

	// 2. 모든 내용 tracing 적용.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// 3.  parent spand를 만들어준다.
	ctx, span := trace.StartSpan(context.Background(), "main")
	defer span.End()

	for i := 0; i < 10; i++ {
		doWork(ctx)
		doWork2(ctx)
	}
}

func doWork(ctx context.Context) {
	//4. child span 을 시작한다.

	_, span := trace.StartSpan(ctx, "functions 1")
	defer span.End()

	fmt.Println("doing busy work")
	time.Sleep(80 * time.Millisecond)
	buf := bytes.NewBuffer([]byte{0xFF, 0x00, 0x00, 0x00})
	num, err := binary.ReadVarint(buf)
	if err != nil {
		// 6. error 로그 설정
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
	}

	// 7. span에 대한 로그 설정
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("bytes to int", num),
	}, "Invoking doWork")
	time.Sleep(20 * time.Millisecond)
}

func doWork2(ctx context.Context) {
	//4. child span 을 시작한다.

	_, span := trace.StartSpan(ctx, "functions 2")
	defer span.End()

	fmt.Println("doing busy work")
	time.Sleep(80 * time.Millisecond)
	buf := bytes.NewBuffer([]byte{0xFF, 0x00, 0x00, 0x00})
	num, err := binary.ReadVarint(buf)
	if err != nil {
		// 6. error 로그 설정
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
	}

	// 7. span에 대한 로그 설정
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("bytes to int", num),
	}, "Invoking doWork")
	time.Sleep(20 * time.Millisecond)
}
