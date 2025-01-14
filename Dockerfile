# Use the official Golang Alpine image
FROM golang:alpine

# Install required packages for building the project
RUN apk update && \
    apk add --no-cache build-base clang linux-headers libbpf-dev llvm-dev

# Set the working directory
WORKDIR /demo

# Copy the project source code to the working directory
COPY . .

# Download the Go module dependencies
RUN go mod download

# Generate eBPF code and other assets as needed
RUN go generate ./...

# Build the Go application
RUN go build -o xdp-demo .

# Set the entry point for the container
ENTRYPOINT ["./xdp-demo"]
