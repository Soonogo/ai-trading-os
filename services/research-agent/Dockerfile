FROM golang:1.25-alpine AS builder
WORKDIR /app/services/research-agent
COPY . /app
RUN go mod download
RUN go build -o /app/bin/research-agent .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/research-agent /app/research-agent
EXPOSE 8083
ENTRYPOINT ["/app/research-agent"]
