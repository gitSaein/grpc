package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	pb "grpc/opencensus/rpc"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	testCount *int
)

func init() {
	testCount = flag.Int("count", 1, "input count")
	flag.Parse()
}

const indent = "  "

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

func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func main() {

	rand.Seed(time.Now().UnixNano())

	// 측정값 기본으로 console 로 출력, 로컬에서 따로 정의 해서 다르게 출력 한다
	view.RegisterExporter(&customMetricsExporter{})

	// // Register the view to collect gRPC client stats.
	// if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
	// 	log.Fatal(err)
	// }
	// 디폴트 값이 아닌 선택해서 출력: RoundtripLatency, CompletedRPC, ServerLatency
	if err := view.Register(ocgrpc.ClientRoundtripLatencyView, ocgrpc.ClientCompletedRPCsView, ocgrpc.ClientServerLatencyView); err != nil {
		log.Fatal(err)
	}

	// Set up a connection to the server with the OpenCensus
	// stats handler to enable stats and tracing.
	conn, err := grpc.Dial(address, grpc.WithStatsHandler(&ocgrpc.ClientHandler{}), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName

	// reporting 은 1초 마다 갱신 한다
	view.SetReportingPeriod(time.Second)

	// 컴맨드에서 입력 받은 값으로 전송 횟수를 조절, 랜덤으로 2가지 rpc 호출하고 갯수 저장
	sayHelloCnt := 0
	capitalizeCnt := 0
	for i := 0; i < *testCount; i++ {
		log.Printf("[%d]---------------------------------------------------------------------------", i+1)

		if randInt(1, 2) == 1 {
			r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
			if err != nil {
				log.Printf("Could not greet: %v", err)
			} else {
				log.Printf("Greeting: %s", r.Message)
			}

			sayHelloCnt++
		} else {
			r, err := c.Capitalize(context.Background(), &pb.Payload{Data: name})
			if err != nil {
				log.Printf("Could not greet: %v", err)
			} else {
				log.Printf("Capitalize: %s", r.Data)
			}
			capitalizeCnt++
		}

		time.Sleep(5 * time.Second)
	}

	log.Printf("sayHelloCnt:%d, capitalizeCnt:%d", sayHelloCnt, capitalizeCnt)
}
