# RPC-bench

RPC-bench is a benchmarking tool designed to compare the performance of gRPC and ttrpc servers. It includes scripts to run benchmarks and generate latency comparison plots.

## Building the Docker Image

To build the Docker image for this project, follow these steps:

1. Clone the repository:
    ```sh
    git clone <repository-url>
    cd rpc-bench
    ```

2. Build the Docker image:
    ```sh
    docker build -t rpc-bench .
    ```

## Running the Benchmark

To run the benchmark from the command line, use the following command:

```sh
# run docker without volume mount
docker run --rm rpc-bench

# run docker with volume mount on output directory
docker run --rm -v $(pwd)/output:/app/output rpc-bench
```
By default, this will run the benchmarks for both gRPC and ttrpc servers and output the results to the `output` directory.
You can modify the `bench.sh` to change the plot generation parameters.
