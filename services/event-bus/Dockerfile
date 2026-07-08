# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app/services/event-bus
COPY . /app
RUN go mod download
RUN go build -o /app/bin/event-bus .

# Runtime stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/event-bus /app/event-bus
EXPOSE 8080
ENTRYPOINT ["/app/event-bus"]
