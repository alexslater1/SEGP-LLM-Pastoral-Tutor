# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o worker-node .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required dependencies for onnxruntime
RUN apk add --no-cache \
  libc6-compat \
  libstdc++ \
  gcompat

# Copy the binary from builder
COPY --from=builder /app/worker-node .

# Copy model directory and ensure it's in the correct location
COPY model/ ./model/

# Copy libonnxruntime and ensure it's in the correct location
COPY libonnxruntime.so.1.20.1 .

# Create a directory for .env file
COPY .env .

# Make the binary executable
RUN chmod +x ./worker-node

# Set LD_LIBRARY_PATH to include current directory
ENV LD_LIBRARY_PATH=/app:$LD_LIBRARY_PATH

ENTRYPOINT ["./worker-node"] 