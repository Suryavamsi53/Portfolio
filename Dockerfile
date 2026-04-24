# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 2: Final lightweight image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Copy static files (HTML, CSS, JS, etc.)
COPY index.html .
COPY login.html .
COPY dashboard.html .
COPY style.css .
COPY script.js .
COPY logo.png .
# Note: messages.json will be created by the app if it doesn't exist,
# but for persistence, you should mount a volume here.
# COPY messages.json . 

# Expose the port the app runs on
EXPOSE 8001

# Command to run the executive
CMD ["./server"]
