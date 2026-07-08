
FROM golang:1.25 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/backtest-engine .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /bin/backtest-engine /backtest-engine
ENTRYPOINT ["/backtest-engine"]
