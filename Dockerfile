# Build stage
FROM golang:1.23-alpine AS builder

# Install git (required for fetching dependencies)
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o object-storage-service main.go

# Final stage: minimal image
FROM scratch

# Copy the binary from the builder
COPY --from=builder /app/object-storage-service /object-storage-service

# Expose the port (default to 8080, but you can override via ENV or args)
EXPOSE 8080

# Command to run the service
ENTRYPOINT ["/object-storage-service"]
