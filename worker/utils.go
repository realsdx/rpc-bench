package worker

import (
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

// set zerlog console logegr
var logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

// Helper function to generate a random string of a given length
func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Print Statistics in a formatted way
func PrintStats(latencies *[]time.Duration, timeElapsed time.Duration) {
	sort.Slice(*latencies, func(i, j int) bool {
		return (*latencies)[i] < (*latencies)[j]
	})

	var sum time.Duration
	for _, lat := range *latencies {
		sum += lat
	}

	numCalls := len(*latencies)
	avg := sum / time.Duration(numCalls)
	median := (*latencies)[numCalls/2]
	min := (*latencies)[0]
	max := (*latencies)[numCalls-1]

	p99 := (*latencies)[int(float64(numCalls)*0.99)]
	p90 := (*latencies)[int(float64(numCalls)*0.90)]

	// calculate request per second
	rps := float64(numCalls) / float64(timeElapsed)

	// Print results
	logger.Info().Msg(" ------------------------ ")
	logger.Info().Msgf("Benchmark Results (%d calls):\n", numCalls)
	logger.Info().Msgf("  Average:   %s", avg)
	logger.Info().Msgf("  Median:    %s", median)
	logger.Info().Msgf("  Minimum:   %s", min)
	logger.Info().Msgf("  Maximum:   %s", max)
	logger.Info().Msgf("  99th %%:    %s", p99)
	logger.Info().Msgf("  90th %%:    %s", p90)
	logger.Info().Msgf("  RPS:       %.2f", rps)
	logger.Info().Msg(" ------------------------ ")

}

// Write all latency to a file in microseconds
func WriteData(latencies *[]time.Duration, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		logger.Fatal().Msgf("Could not create file: %v", err)
	}
	defer file.Close()

	for _, lat := range *latencies {
		// write the latency in microseconds
		file.WriteString(strconv.FormatInt(lat.Microseconds(), 10) + "\n")
	}
	logger.Info().Msgf("Latency data written to %s", filename)
}
