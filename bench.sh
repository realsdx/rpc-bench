#!/bin/bash

set -e

# Function to display help text
show_help() {
    echo "Usage: $0 -c <concurrency> -t <time>"
    echo ""
    echo "Options:"
    echo "  -c, --concurrency   Number of concurrent threads to use"
    echo "  -t, --time          Time duration for the benchmark in seconds"
    echo "  -h, --help          Display this help message"
}

# Default values
concurrency=""
time=""

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -c|--concurrency) concurrency="$2"; shift ;;
        -t|--time) time="$2"; shift ;;
        -h|--help) show_help; exit 0 ;;
        *) echo "Unknown parameter passed: $1"; show_help; exit 1 ;;
    esac
    shift
done

# Check if required arguments are provided
if [[ -z "$concurrency" || -z "$time" ]]; then
    show_help
    exit 1
fi

# Run the grpc & ttrpc servers , output is redirected to /dev/null
echo "Starting grpc test server at port 50051"
go run testservers/grpc/main.go  -port=50051 > /dev/null 2>&1 &

echo "Starting ttrpc test server at port 50052"
go run testservers/ttrpc/main.go -port=50052 > /dev/null 2>&1 &

sleep 2

# Run tbench binary
echo ""
echo "Benchmarking ttrpc server..."
echo "Concurrency: $concurrency | Time: $time seconds"
./tbench -c="$concurrency" -t="$time" 127.0.0.1:50052

# Run gbench binary
echo ""
echo "Benchmarking grpc server..."
echo "Concurrency: $concurrency | Time: $time seconds"
./gbench -c="$concurrency" -t="$time" 127.0.0.1:50051

echo "Benchmarking completed"

# Kill the servers
echo "Killing test servers..."
lsof -i :50051 | awk 'NR!=1 {print $2}' | xargs kill -9
lsof -i :50052 | awk 'NR!=1 {print $2}' | xargs kill -9

# Run the python script in scripts/plot.py to generate the plot, takes first 10000 latency values
# Change the sigma values in script to smooth the plot as required
python3 scripts/plot.py 5000

echo "Done!"