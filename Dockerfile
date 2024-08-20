# Use the official Golang image as the base image
FROM golang:1.23

# Install the required dependencies
RUN apt-get update && apt-get install -y \
    lsof \
    python3-pip

ENV PIP_BREAK_SYSTEM_PACKAGES=1
RUN pip3 install --upgrade pip && \
    pip3 install matplotlib numpy scipy seaborn

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the project dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the Go applications
RUN go build -o tbench ./cmd/tbench
RUN go build -o gbench ./cmd/gbench

# Set the command to run the executables
CMD ["./bench.sh", "-c", "2", "-t", "5"]