package main

import (
	"flag"
	"time"

	"github.com/realsdx/ttrpc-bench/worker"
)

func main() {
	// parse flags for concurrenty and time of run
	// flags
	concurrency := flag.Int("c", 1, "Number of concurrent requests")
	duration := flag.Int("t", 5, "How long to run this test. In seconds")
	flag.Parse()

	host := flag.Arg(0)

	runner := worker.TTRPCRunner{
		ServerAddress: host,
		Concurrency:   *concurrency,
		Duration:      time.Duration(*duration),
	}
	runner.Run()
}
