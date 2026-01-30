# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/controller ./cmd/controller && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/transfer ./cmd/transfer && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/agent ./cmd/agent

# Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/controller /app/controller
COPY --from=builder /app/transfer /app/transfer
COPY --from=builder /app/agent /app/agent

# No default CMD; Helm sets command per Deployment (controller / transfer / agent)
