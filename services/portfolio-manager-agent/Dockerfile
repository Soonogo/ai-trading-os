FROM golang:1.25-alpine AS builder
WORKDIR /app/services/portfolio-manager-agent
COPY . /app
RUN go mod download
RUN go build -o /app/bin/portfolio-manager-agent .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/portfolio-manager-agent /app/portfolio-manager-agent
EXPOSE 8084
ENTRYPOINT ["/app/portfolio-manager-agent"]
