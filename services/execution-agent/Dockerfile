FROM golang:1.25-alpine AS builder
WORKDIR /app/services/execution-agent
COPY . /app
RUN go mod download
RUN go build -o /app/bin/execution-agent .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/execution-agent /app/execution-agent
ENTRYPOINT ["/app/execution-agent"]
