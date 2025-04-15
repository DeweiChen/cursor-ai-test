# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o chatroom-app

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/chatroom-app .
# Copy the client.html file
COPY --from=builder /app/client.html .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./chatroom-app"] 