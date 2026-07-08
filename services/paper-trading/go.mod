module github.com/Soonogo/ai-trading-os/services/paper-trading

go 1.25.0

require github.com/Soonogo/ai-trading-os/services/common v0.0.0

require (
	github.com/google/uuid v1.3.1 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/nats-io/nats.go v1.51.0 // indirect
	github.com/nats-io/nkeys v0.4.16 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
)

replace github.com/Soonogo/ai-trading-os/services/common => ../common
