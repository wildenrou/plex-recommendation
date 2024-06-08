# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22.3-alpine@sha256:b8ded51bad03238f67994d0a6b88680609b392db04312f60c23358cc878d4902 AS builder

WORKDIR /app

# Copy Go Modules and install dependencies 
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Run tests
RUN go test ./...

# Build the Go binary statically 
RUN CGO_ENABLED=0 GOOS=linux go build -o recommendations ./cmd/main.go

# Final stage
FROM alpine:3.20.0@sha256:77726ef6b57ddf65bb551896826ec38bc3e53f75cdde31354fbffb4f25238ebd

# Create apprunner user and group
RUN addgroup -S apprunner && adduser -S apprunner -G apprunner

WORKDIR /app

# Only copy the necessary binary and environment file
COPY --from=builder /app/recommendations /app/recommendations
COPY --from=builder /app/.env /app/.env

# Change ownership to apprunner user
RUN chown -R apprunner:apprunner /app

# Expose API port
EXPOSE 8090

# Command to run the application as apprunner user
USER apprunner
CMD ["/app/recommendations"]
