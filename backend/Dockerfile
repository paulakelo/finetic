# STAGE 1: Build the Go binary
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

# Download dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the self-contained Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /finetic-api ./cmd/api

# STAGE 2: Lightweight runtime
FROM alpine:latest

WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /finetic-api .

# Expose the port the application runs on
EXPOSE 8080

# Run the Go binary
CMD ["./finetic-api"]