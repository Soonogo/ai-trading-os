# AITOS Architecture

## Overview

AITOS is an event-driven, multi-agent investment platform. Every AI decision is recorded with reasoning, confidence, and traceability.

## Architectural principles

1. **Event-sourced cognition** — agents publish decisions as events; the system replays decisions for audits and backtests.
2. **Multi-model debate** — the Model Router calls several LLMs and aggregates confidence using a weighted scoring layer.
3. **Explainability first** — every signal includes a `ReasoningTimeline` with evidence and source IDs.
4. **Domain-driven design (DDD)** — each Go microservice owns a bounded context: aggregates, repositories, domain events.
5. **CQRS-lite** — command-side writes via agents; read-side projections feed dashboards and reports.

## Runtime topology

```
┌──────────────┐     WebSocket/NATS      ┌────────────────┐
│   Next.js    │<───────────────────────>│   API Gateway  │
│   (apps/web) │                         │   (event-bus)  │
└──────────────┘                         └────────────────┘
                                                │
                          ┌─────────────────────┼──────────────────────┐
                          │                     │                      │
              Research    │  Macro   News      Technical  Fundamental  Risk
              Agent       │  Agent   Agent     Agent      Agent        Agent
                          │                                      │
                          └──────────────────────┬───────────────┘
                                                 │
                              Portfolio Manager Agent
                                                 │
                                    Execution Agent
                                                 │
                                     Exchanges / Brokers
```

## Agent responsibilities

| Agent | Responsibility | Writes events | Reads events |
|-------|----------------|---------------|--------------|
| Research | Multi-asset screen, RAG over SEC/whitepapers/filings | `ResearchSignal` | `MarketData`, `NewsAnalyzed` |
| Macro | Economic regime, Fed, yields, global flows | `MacroSignal` | `EconomicCalendar`, `MarketData` |
| News | NLP over headlines, sentiment, surprise scoring | `NewsSignal` | `RawNews` |
| Technical | Pattern, indicator, volume, market structure | `TechnicalSignal` | `MarketData`, `OrderBook` |
| Fundamental | Earnings, valuation, quality metrics | `FundamentalSignal` | `FundamentalsUpdated` |
| Risk | VaR, position heat, correlation, drawdown | `RiskSignal` | All signal events |
| Portfolio | Optimization, target weights, allocation drift | `PortfolioDecision` | All signal + risk events |
| Execution | Order slicing, venue selection, smart routing | `OrderSubmitted`, `OrderFilled` | `PortfolioDecision` |
| Reflection | Reviews past decisions, audits mistakes | `ReflectionReport` | All decisions |
| Report | Daily summaries, P&L attribution | `DailyReport` | All events |

## Event schema

Every event has:

- `event_id` (UUID v7)
- `trace_id` (links decisions across services)
- `version` (schema version)
- `timestamp` (UTC)
- `source_agent`
- `reasoning` (JSON pointer to `ReasoningTimeline`)

See `packages/event-schemas` for canonical Go/JSON definitions.

## Data stores

| Store | Use |
|-------|-----|
| PostgreSQL | OLTP: accounts, portfolios, orders, agents, reasoning records |
| Redis | Caches, rate limits, session state, real-time leaderboards |
| ClickHouse | Time-series market data, P&L, backtests, analytics |
| Qdrant | Vector memory and RAG documents |
| NATS JetStream | Event bus and persistent queues |

## Security & operations

- All service-to-service calls require mTLS or NATS credentials.
- Live trading requires explicit `PAPER=false` and broker API keys.
- Sensitive keys live in `/run/repo_secrets/{repo}/.env.secrets` and are never committed.
