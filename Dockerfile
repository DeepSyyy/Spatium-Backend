# ---------- STAGE 1: Build ----------
FROM golang:1.25.1-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies for building (optional if you use cgo)
RUN apk add --no-cache git

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the app binary
RUN go build -o main .

# ---------- STAGE 2: Runtime ----------
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose API port
EXPOSE 4419

# Environment variables (Railway will override this if needed)
ENV PORT=4419

# Run the app
CMD ["./main", "serve"]