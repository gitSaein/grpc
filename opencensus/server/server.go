package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	pb "grpc/opencensus/rpc"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/runmetrics"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":50051"

const indent = "  "

type server struct {
	pb.UnimplementedGreeterServer
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	// 1 ~ 3 랜덤
	__delay := randInt(1, 3)

	log.Printf("[Request] SayHello:%s(%d sec)", in.Name, __delay)

	// latency 테스트를 위해 1 ~ 3초 사이로 랜덤으로 delay 를 시킨다
	time.Sleep(time.Duration(__delay) * time.Second)

	return &pb.HelloReply{
		Message: "Hello " + in.Name,
	}, nil
}

func (s *server) Capitalize(ctx context.Context, in *pb.Payload) (*pb.Payload, error) {

	// 1 ~ 3 랜덤
	__delay := randInt(1, 3)

	log.Printf("[Request] Capitalize:%s(%d sec)", in.Data, __delay)

	// latency 테스트를 위해 1 ~ 3초 사이로 랜덤으로 delay 를 시킨다
	time.Sleep(time.Duration(__delay) * time.Second)

	return &pb.Payload{
		Data: strings.ToUpper(in.Data),
	}, nil
}

type customMetricsExporter struct{}

func (ce *customMetricsExporter) ExportView(vd *view.Data) {

	if len(vd.Rows) == 0 {
		return
	}

	log.Printf("--------------------------------")

	for _, row := range vd.Rows {
		fmt.Printf("%v %-45s", vd.End.Format("15:04:05"), vd.View.Name)

		switch v := row.Data.(type) {
		case *view.DistributionData:
			fmt.Printf("distribution: min=%.1f max=%.1f mean=%.1f", v.Min, v.Max, v.Mean)
		case *view.CountData:
			fmt.Printf("count:        value=%v", v.Value)
		case *view.SumData:
			fmt.Printf("sum:          value=%v", v.Value)
		case *view.LastValueData:
			fmt.Printf("last:         value=%v", v.Value)
		}
		fmt.Println()

		for _, tag := range row.Tags {
			fmt.Printf("%v- %v=%v\n", indent, tag.Key.Name(), tag.Value)
		}
	}
	log.Printf("--------------------------------")
}

type customTraceExporter struct{}

func (e *customTraceExporter) ExportSpan(vd *trace.SpanData) {
	var (
		traceID      = hex.EncodeToString(vd.SpanContext.TraceID[:])
		spanID       = hex.EncodeToString(vd.SpanContext.SpanID[:])
		parentSpanID = hex.EncodeToString(vd.ParentSpanID[:])
	)
	var reZero = regexp.MustCompile(`^0+$`)

	fmt.Println()
	fmt.Println("#----------------------------------------------")
	fmt.Println()
	fmt.Println("TraceID:     ", traceID)
	fmt.Println("SpanID:      ", spanID)
	if !reZero.MatchString(parentSpanID) {
		fmt.Println("ParentSpanID:", parentSpanID)
	}

	fmt.Println()
	fmt.Printf("Span:    %v\n", vd.Name)
	fmt.Printf("Status:  %v [%v]\n", vd.Status.Message, vd.Status.Code)
	fmt.Printf("Elapsed: %v\n", vd.EndTime.Sub(vd.StartTime).Round(time.Millisecond))

	if len(vd.Annotations) > 0 {
		fmt.Println()
		fmt.Println("Annotations:")
		for _, item := range vd.Annotations {
			fmt.Print(indent, item.Message)
			for k, v := range item.Attributes {
				fmt.Printf(" %v=%v", k, v)
			}
			fmt.Println()
		}
	}

	if len(vd.Attributes) > 0 {
		fmt.Println()
		fmt.Println("Attributes:")
		for k, v := range vd.Attributes {
			fmt.Printf("%v- %v=%v\n", indent, k, v)
		}
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	err := runmetrics.Enable(runmetrics.RunMetricOptions{
		EnableCPU:    true,
		EnableMemory: true,
		Prefix:       "mayapp/",
	})
	if err != nil {
		log.Fatal(err)
	}

	exp_pro, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "prometheus_grpc", // prometheus에서 읽어 들일 카테고리 이름,  prometheus.yaml 에 셋팅값으로 들어간다
	})
	if err != nil {
		log.Fatal(err)
	}

	// 지표가 수집될 웹서버를 생성
	// 측정값 웹으로 보내주고 prometheus가 metrics 을 가져가게 된다
	// localhost:9888/metrics --> 들어가면 지표값을 확인 할 수 있다``
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", exp_pro) // prometheus와 연결
		log.Fatal(http.ListenAndServe(":9888", exp_pro))
	}()

	// 측정값 기본으로 console 로 출력, 로컬에서 따로 정의해서 사용자 원하는 출력으로 바꿔서 사용
	view.RegisterExporter(&customMetricsExporter{})
	trace.RegisterExporter(&customTraceExporter{})

	// Register the views to collect server request count.
	// if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
	// 	log.Fatal(err)
	// }
	// 디폴트 값이 아닌 선택해서 출력: ServerLatency, ServerCompletedRPCs(Count)
	if err := view.Register(ocgrpc.ServerLatencyView, ocgrpc.ServerCompletedRPCsView); err != nil {
		log.Fatal(err)
	}

	// reporting 은 5초 마다 갱신 한다
	view.SetReportingPeriod(5 * time.Second)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// opencensus 에 핸들러 등록
	s := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s) // grpcurl 명령을 사용하게 하기 위해

	log.Printf("start gRPC server on %s port", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
