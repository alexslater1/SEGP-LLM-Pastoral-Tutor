# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o worker-node .

# Final stage
FROM ubuntu:latest

WORKDIR /app

# Install required dependencies
RUN apt-get update && apt-get install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder
COPY --from=builder /app/worker-node .

# Copy model directory and ensure it's in the correct location
COPY model/ ./model/

# Copy libonnxruntime and ensure it's in the correct location
COPY libonnxruntime.so.1.20.1 .

# Copy .env file
COPY .env .

# Make the binary executable
RUN chmod +x ./worker-node

# Set LD_LIBRARY_PATH to include current directory
ENV LD_LIBRARY_PATH=/app:$LD_LIBRARY_PATH

ENTRYPOINT ["./worker-node"] 