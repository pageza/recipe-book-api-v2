#Dockerfile

# Use Ubuntu 22.04 as the base image
FROM ubuntu:22.04

# Set noninteractive mode for apt-get to avoid prompts during build
ARG DEBIAN_FRONTEND=noninteractive

# Install required packages
RUN apt-get update && apt-get install -y \
    curl \
    git \
    ca-certificates \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.24.0 (the latest stable version per your instruction)
ENV GO_VERSION=1.24.0
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -o go.tar.gz \
    && tar -C /usr/local -xzf go.tar.gz \
    && rm go.tar.gz

# Update PATH to include Go binaries
ENV PATH="/usr/local/go/bin:${PATH}"

# Set the working directory
WORKDIR /workspace

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o recipe-app ./cmd/grpc-server/main.go

# Expose the API port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/workspace/recipe-app"]
