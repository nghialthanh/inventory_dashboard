# Inventory Dashboard

Dealership **inventory dashboard**: filter vehicles, identify **aging stock** (in stock more than 90 days), and record a **proposed action**.

**Assessment choice:** this repository’s graded service layer is the **backend** (REST API + PostgreSQL). A React UI is included as a live client for demos; you can exercise the API with the test suite or [cURL examples](docs/api-curl-examples.md) instead.

**System design (Part 1):** [docs/SYSTEM_DESIGN.md](docs/SYSTEM_DESIGN.md)

## Stack

| Layer | Choice |
|-------|--------|
| Backend (assessed) | Golang modular monolith, REST |
| Logging / metrics | zap + `X-Trace-ID`; Prometheus `GET /metrics` |
| Database | PostgreSQL primary + streaming replica (catalog list reads) |
| Cache | Redis cache-aside for catalog make/model lists ([ADR-0014](docs/adr/0014-catalog-cache-and-read-replica.md)) |
| Frontend (bonus client) | React + TypeScript |
| Local run | Docker Compose — API/UI `:8080`, Postgres `:5432`, replica `:5433`, Redis `:6379` |

## Where to read

| Doc | Purpose |
|-----|---------|
| [docs/SYSTEM_DESIGN.md](docs/SYSTEM_DESIGN.md) | Architecture, data flow, tech justifications, observability, GenAI in design |
| [GLOSSARY.md](GLOSSARY.md) | Domain language |
| [docs/adr/](docs/adr/) | Hard decisions (0001–0014). The [index](docs/adr/README.md) pairs each decision with the issue that accepted it |
| [docs/architecture/](docs/architecture/) | API outline, schema, directory layout |
| [docs/api-curl-examples.md](docs/api-curl-examples.md) | Client stub via cURL |
| [.scratch/inventory-dashboard/spec.md](.scratch/inventory-dashboard/spec.md) | Scenario behavior |

## Build and run

Docker Compose serves the API and the built screen together on port 8080, with Postgres primary (`5432`), replica (`5433`), and Redis (`6379`). There is no sign-in.

Catalog addresses are integers ([ADR-0011](docs/adr/0011-integer-addresses.md)); vehicle addresses are snowflakes ([ADR-0013](docs/adr/0013-snowflake-vehicle-addresses.md)). Drop Compose volumes before starting after an identity/schema or replication-init change:

```bash
docker compose down -v
docker compose up --build
```

Open http://localhost:8080. Liveness: `GET /health`. Metrics: `GET /metrics` (outside `/api/v1`).

Optional UI-only dev (API still from Compose):

```bash
cd frontend && npm install && npm run dev
```

Vite proxies `/api` and `/health` to `http://localhost:8080`. See [frontend/README.md](frontend/README.md).

Optional bulk demo data (idempotent; ≥510 vehicles per monthly partition 2026-01…2027-12):

```bash
docker compose exec -T postgres psql -U inventory -d inventory < scripts/seed_bulk_vehicles.sql
```

## Test

Core business rules (aging boundaries, VIN uniqueness, terminal stock-out, proposed action) are covered by HTTP API tests against PostgreSQL with a fixed UTC clock:

```bash
cd backend && go test ./...
```

## Status

Implemented: health, catalog, vehicles (record / correct / list filters / stock-out / proposed action), Compose one-origin deploy, and the inventory React screen as a bonus client.

---

## AI Collaboration Narrative

**Full step-by-step (recommended for submission / interview):** [docs/AI_COLLABORATION.md](docs/AI_COLLABORATION.md)

Design-phase summary: [docs/SYSTEM_DESIGN.md §7](docs/SYSTEM_DESIGN.md#7-genai-assistance-in-the-design-phase).

### Strategy for guiding the AI

1. **Lock the vocabulary first** — prompts and reviews use [GLOSSARY.md](GLOSSARY.md) terms only (`aging stock`, `stocked out`, `proposed action`).
2. **Architecture before code** — ADRs and `docs/architecture/*`; handlers bind/map errors; logic owns rules; repositories own SQL.
3. **Small, testable slices** — one command or invariant at a time (aging 90 vs 91, terminal stock-out, proposed action only when aging).
4. **Contract as source of truth** — [api-outline.md](docs/architecture/api-outline.md); UI shaped to DTOs.

### Verifying and refining output

1. `cd backend && go test ./...` after substantive backend changes.
2. Assert envelope, `error_code`, and derived fields (`stock_state`, `inventory_age`, `is_aging`).
3. Fixed UTC clock for aging boundaries.
4. Update the owning doc in the same change when behavior moves.
5. Reject ADR drift (free-text make/model, stock-state column, stock-out on PATCH).

### Ensuring final code quality

| Gate | Practice |
|------|----------|
| Domain | Tests for age 0 / 90 / 91, stocked-out not aging, duplicate VIN, command isolation |
| Structure | Composition root in `server.go`; no inventory rules in handlers |
| Operability | Compose (primary/replica/Redis), `/health`, `/metrics`, graceful shutdown, zap + `trace_id` |
| Review | Human review of GenAI diffs against glossary, ADRs, and the spec |

GenAI accelerated drafts and scaffolding. **I owned architecture and acceptance**, using ADRs and automated tests as the gate.
