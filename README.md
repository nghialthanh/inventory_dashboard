# Inventory Dashboard

Dealership **inventory dashboard** assessment: filter vehicles, identify **aging stock** (in stock more than 90 days), and record a **proposed action**.

## Stack

| Layer | Choice |
|-------|--------|
| Backend | Golang (modular monolith, REST) |
| Frontend | React + TypeScript |
| Database | PostgreSQL |
| Images | MinIO (private; streamed by API) |
| Local run | Docker Compose — one origin |

## Domain language & decisions

- **Glossary:** [GLOSSARY.md](GLOSSARY.md) — use these words in code, API, UI, and docs.
- **ADRs:** [docs/adr/](docs/adr/) — retain stocked-out vehicles, catalog make/model, MinIO primary image, terminal stocked-out, derived stock state, UUID identity, separate commands, API-streamed images, one-origin Compose.
- **Architecture:** [docs/architecture/README.md](docs/architecture/README.md)

## Cursor rules

Project rules mirror the Homestay ops layout, adapted to this domain:

| Rule | Role |
|------|------|
| `.cursor/rules/architecture.mdc` | Always-on modular monolith + glossary/ADR discipline |
| `.cursor/rules/frontend-integration.mdc` | Always-on FE presentation SoT |
| `.cursor/rules/backend.mdc` | Go handlers / logic / repositories |
| `.cursor/rules/frontend.mdc` | React hooks / adapters |
| `.cursor/rules/api.mdc` | `/api/v1` envelope + commands |
| `.cursor/rules/database.mdc` | Postgres schema conventions |
| `.cursor/rules/ui.mdc` | Dashboard UI patterns |
| `.cursor/rules/testing.mdc` | Aging / VIN / stock-out tests |
| `.cursor/rules/documentation.mdc` | Keep docs in sync |

## Architecture docs

| Doc | Contents |
|-----|----------|
| [docs/architecture/README.md](docs/architecture/README.md) | System context, domain, workflows |
| [directory-structure.md](docs/architecture/directory-structure.md) | Repo layout |
| [database-schema.md](docs/architecture/database-schema.md) | Tables & indexes |
| [api-outline.md](docs/architecture/api-outline.md) | REST contract |
| [frontend-integration.md](docs/architecture/frontend-integration.md) | FE-first integration |
| [frontend-backend-mapping.md](docs/architecture/frontend-backend-mapping.md) | Page ↔ API ↔ hooks |
| [inventory-frontend.md](docs/architecture/inventory-frontend.md) | UI stack notes |

## Brief

See [initial-prompt.md](initial-prompt.md).

## Status

Design artifacts (glossary, ADRs, architecture, Cursor rules) are in place. Application code (`backend/`, `frontend/`, Compose) is not scaffolded yet — implement incrementally against the architecture docs.
