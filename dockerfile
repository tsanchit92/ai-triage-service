# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go build -o server ./cmd/server

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/internal/migrations ./internal/migrations
COPY conf.env .
EXPOSE 8080
ENTRYPOINT ["./server"]
