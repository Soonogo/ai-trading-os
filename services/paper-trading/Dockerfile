
FROM golang:1.25 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/paper-trading .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /bin/paper-trading /paper-trading
ENTRYPOINT ["/paper-trading"]
