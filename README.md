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
│   └── backtest-engine
├── packages/                # Shared code
│   ├── design-system        # UI primitives & Tailwind preset
│   ├── event-schemas        # Shared Go event definitions
│   └── tsconfig             # Shared TS configs
├── docs/                    # Architecture ADRs and specs
├── docker-compose.yml
└── .github/workflows/ci.yml
```

## Quick start

```bash
cp .env.example .env
docker compose up -d db redis clickhouse qdrant nats
pnpm install
pnpm --filter @aitos/web dev
# In another terminal
cd services/event-bus && go run .
```

## License

Proprietary — see `LICENSE`.
