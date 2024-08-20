package worker

import (
	"context"
	"net"
	"time"

	"github.com/containerd/ttrpc"
	pb "github.com/realsdx/ttrpc-bench/testproto/pbttrpc"
)

type TTRPCRequester struct {
	conn   net.Conn
	client pb.GreeterService
}

func NewTTRPCRequester(serverHost string) *TTRPCRequester {
	// create a new connection
	conn, err := net.Dial("tcp", serverHost)
	if err != nil {
		logger.Fatal().Msgf("did not connect: %v", err)
	}

	client := pb.NewGreeterClient(ttrpc.NewClient(conn))

	return &TTRPCRequester{conn: conn, client: client}
}

func (req *TTRPCRequester) Close() {
	req.conn.Close()
}

// make ttrpc request to server
func (req *TTRPCRequester) CallRPC() time.Duration {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t1 := time.Now()
	_, err := req.client.SayHello(ctx, &pb.HelloRequest{Name: GetRandomString(8)})
	timeElapsed := time.Since(t1)
	if err != nil {
		logger.Fatal().Msgf("could not greet: %v", err)
	}
	// logger.Info().Msgf("Greeting: %s, latency: %s", r.GetMessage(), timeElapsed)

	return timeElapsed
}

type TTRPCRunner struct {
	ServerAddress string
	Concurrency   int
	Duration      time.Duration
}

func (r *TTRPCRunner) Run() {
	logger.Info().Msg("Starting ttRPC Benchmark")

	// Set timeout context
	timeout := r.Duration * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stopLoop := make(chan struct{}, 1)
	semaphore := make(chan struct{}, r.Concurrency)

	// initialize requester
	req := NewTTRPCRequester(r.ServerAddress)

	// warm up
	for i := 0; i < 200; i++ {
		req.CallRPC()
	}
	logger.Info().Msg("Warmup Complete")

	var latencies []time.Duration
	var numCalls int = 0

mainLoop:
	for {
		select {
		case <-stopLoop:
			// logger.Info().Msg("StopLoop Called")
			break mainLoop
		default:
			semaphore <- struct{}{} // acquire semaphore

			go func(server_host string) {
				// defer semaphore release
				defer func() { <-semaphore }()

				select {
				case <-ctx.Done():
					// logger.Info().Msgf("Timeout occured after %d seconds", r.Duration)
					stopLoop <- struct{}{}
				default:
					// make ttrpc request
					latency := req.CallRPC()

					latencies = append(latencies, latency)
					numCalls += 1
				}
			}(r.ServerAddress)
		}
	}

	WriteData(&latencies, "output/ttrpc_latencies.txt")
	PrintStats(&latencies, r.Duration)
}
