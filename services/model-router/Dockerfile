# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app/services/model-router
COPY . /app
RUN go mod download
RUN go build -o /app/bin/model-router .

# Runtime stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/model-router /app/model-router
EXPOSE 8082
ENTRYPOINT ["/app/model-router"]
