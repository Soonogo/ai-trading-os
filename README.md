# AITOS — AI Trading Operating System

An autonomous, multi-agent investment platform.

## Product feel

Notion + Linear + TradingView + Perplexity + Bloomberg Terminal. Dark, glassmorphic, minimal, AI-first.

## Tech stack

- **Frontend:** React 19, Next.js 15, TypeScript, Tailwind CSS, shadcn/ui, Motion, TradingView chart, TanStack Query, Zustand
- **Backend:** Go, Gin, PostgreSQL, Redis, ClickHouse, Qdrant, NATS (JetStream)
- **AI:** multi-agent pipeline with model router, RAG, prompt versioning, confidence scoring, and explainability
- **DevOps:** Docker, Docker Compose, GitHub Actions, Vitest, Playwright

## Monorepo layout

```
.
├── apps/web                 # Next.js 15 dashboard
├── services/                # Go microservices (DDD, one module per service)
│   ├── event-bus
│   ├── model-router
│   ├── research-agent
│   ├── macro-agent
│   ├── news-agent
│   ├── technical-analysis-agent
│   ├── fundamental-agent
│   ├── risk-agent
│   ├── execution-agent
│   ├── portfolio-manager-agent
│   ├── reflection-agent
│   ├── report-agent
│   ├── backtest-engine      # Historical strategy simulation
│   └── paper-trading        # Live paper-broker matching
├── packages/                # Shared code
│   ├── design-system        # UI primitives & Tailwind preset
│   ├── event-schemas        # Shared Go event definitions
│   └── tsconfig             # Shared TS configs
├── services/common          # Shared Go domain types and simulation library
├── docs/                    # Architecture ADRs and specs
├── docker-compose.yml
└── .github/workflows/ci.yml
```

## Quick start

```bash
cp .env.example .env
# 1. Start data stores and message bus
docker compose up -d db redis clickhouse qdrant nats
# 2. Install JS/TS dependencies
pnpm install
# 3. Run the Next.js dashboard (http://localhost:3000)
pnpm --filter @aitos/web dev
```

The frontend talks to the Go microservices over REST and WebSocket; agents talk to each other over NATS JetStream.

## Running the platform

### 1. Configure the environment

Copy the example file and fill in any API keys you need:

```bash
cp .env.example .env
```

`.env` is ignored by git. For local development the only required entries are the infrastructure URLs; AI provider keys and broker keys are optional unless you want live model calls or broker data.

### 2. Start infrastructure

```bash
pnpm infra:up
# or directly:
docker compose up -d db redis clickhouse qdrant nats
```

This starts PostgreSQL, Redis, ClickHouse, Qdrant, and NATS with JetStream.

### 3. Install dependencies

```bash
pnpm install
```

Also make sure Go 1.25+ is on your `PATH` (`/usr/local/go/bin` by default in this repo).

### 4. Run the frontend

```bash
pnpm --filter @aitos/web dev
```

Open `http://localhost:3000`.

### 5. Run the API gateway / event bus

The event bus is the NATS-to-WebSocket/REST adapter that the dashboard connects to:

```bash
cd services/event-bus
go run .
```

It listens on `http://localhost:8080` and exposes `/health`, `/ws`, `/v1/events`, and `/v1/prompts`. OpenAPI docs are served at `/docs/openapi.yaml` and Swagger UI at `/swagger`.

### 6. Run AI and trading agents

Each agent is a standalone Go module. You can run them in separate terminals, or start them all with Docker Compose (see below). Typical local flow:

```bash
# Core AI pipeline
cd services/model-router && go run .
cd services/research-agent && go run .
cd services/macro-agent && go run .
cd services/news-agent && go run .
cd services/technical-analysis-agent && go run .
cd services/fundamental-agent && go run .
cd services/risk-agent && go run .

# Portfolio & execution
cd services/portfolio-manager-agent && go run .
cd services/execution-agent && go run .

# Simulation services
cd services/backtest-engine && go run .
cd services/paper-trading && go run .
```

Most agents require only `NATS_URL`. The research-agent also uses `QDRANT_URL` and an `OPENAI_API_KEY` for RAG. The execution-agent uses `PAPER_TRADING=true` (default) to simulate fills instead of sending live orders.

### 7. Run everything in Docker Compose

To start the whole platform with a single command:

```bash
# Optional: build all images first
docker compose build

# Start the whole stack
docker compose up -d
```

This brings up infrastructure plus `event-bus`, `model-router`, `research-agent`, `portfolio-manager-agent`, `execution-agent`, `backtest-engine`, and `paper-trading`. The frontend is not containerized in the default compose file; run it locally with `pnpm --filter @aitos/web dev` as shown above.

## Backtesting and paper trading quick example

The simulation package is in `services/common/sim`. `backtest-engine` consumes `backtest.request` events and emits `backtest.result` events. `paper-trading` consumes `paper.order` and `market.data` events and emits `paper.fill` events.

To run a quick backtest manually:

```bash
cd services/backtest-engine
go run .
```

Then publish a `backtest.request` event via NATS with a payload containing `symbol`, `candles`, `strategy` (`buy-and-hold` or `sma-cross`), and `config`. The result will be published as `backtest.result`.

To start the paper broker:

```bash
cd services/paper-trading
go run .
```

## Running tests

### Web tests

```bash
pnpm typecheck:web
pnpm lint
pnpm test:unit
pnpm test:e2e
```

`test:e2e` uses Playwright and will start the dev server automatically. Run `pnpm --filter @aitos/web exec playwright install` first if browsers are not installed.

### Go tests

```bash
pnpm go:build
pnpm go:test
```

Or run per module:

```bash
cd services/<service-name>
go test ./...
```

## License

Proprietary — see `LICENSE`.
