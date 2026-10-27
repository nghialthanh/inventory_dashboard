# Directory Structure

Planned layout for the inventory dashboard modular monolith. Adjust only with a matching update to this doc and Cursor rules.

## Repository root

```
inventory_dashboard/
├── README.md
├── GLOSSARY.md
├── initial-prompt.md
├── docs/
│   ├── adr/                      # Hard decisions (0001–0014). README.md pairs each with its issue
│   └── architecture/             # Living contract
├── .cursor/rules/                # Agent coding rules
├── backend/                      # Go API (+ serves static UI in Compose)
├── frontend/                     # React + TypeScript SPA
├── docker/                       # Postgres primary/replica init scripts
└── docker-compose.yml            # Postgres primary + replica, Redis, API
```

---

## Backend (`backend/`)

Go module: set in `go.mod` (e.g. `inventory_dashboard`).

Adapters, logic, and repositories. `internal/services` holds the interfaces between them.

| Layer | Where | Rule |
|-------|--------|------|
| Adapter | `api/` for HTTP. Events and third-party clients belong here too | Call `internal/services` interfaces. Do not import `logic` or `repositories` |
| Services | `internal/services` | Declare the feature interfaces. No function bodies and no SQL |
| Logic | `internal/logic` | Implement those interfaces. Feature rules live here |
| Repository | `internal/repositories` | Feature queries and commands against PostgreSQL |

`internal/app/server.go` is the composition root: it builds repositories, builds logic, and hands the logic to the adapter.

```
backend/
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
├── config/
│   ├── app.development.yaml
│   └── example.yaml
├── api/                                 # HTTP adapter
│   ├── handlers/
│   │   ├── catalog.go
│   │   ├── vehicle.go
│   │   └── health.go
│   └── middlewares/
│       ├── trace.go                     # X-Trace-ID → context for all layers
│       ├── access_log.go                # zap request log with trace_id
│       └── …                            # auth only if assessment requires it
├── internal/
│   ├── app/
│   │   ├── server.go                    # wire repositories → logic → adapter
│   │   └── routes/
│   │       ├── router.go
│   │       └── v1/
│   │           └── routes.go
│   ├── common/
│   │   ├── constant/                    # proposed_action enums, aging threshold
│   │   ├── errors/
│   │   └── snowflake/                   # mint and decode vehicle addresses (ADR-0013)
│   ├── dto/
│   ├── models/
│   │   ├── catalog.go
│   │   └── vehicle.go
│   ├── services/                        # feature interfaces only
│   │   ├── health.go
│   │   ├── catalog.go
│   │   └── vehicles.go
│   ├── logic/                           # implements services
│   │   ├── address/                     # parse a catalog or vehicle address
│   │   ├── catalog/                     # list + Redis cache-aside (ADR-0014)
│   │   ├── vehicle/                     # stock-out, patch rules, VIN
│   │   ├── aging/                       # inventory age + is_aging
│   │   └── health/
│   ├── repositories/                    # SQL queries and commands
│   │   ├── catalog.go                   # list on read pool; ModelMakeID on primary
│   │   └── vehicle.go
├── pkg/
│   ├── cache/                           # Redis client for catalog lists
│   ├── config/
│   ├── database/
│   ├── logger/
│   └── metrics/                         # Prometheus counters/histograms
└── …_test.go                            # beside packages
```

Static UI: the image build copies the `frontend` dist into the API image. Echo serves it beside `/api` ([ADR-0009](../adr/0009-one-origin-under-api.md)).

---

## Frontend (`frontend/`)

```
frontend/
├── package.json
├── src/
│   ├── main.tsx
│   ├── routes/                          # or pages/ — pick one router style
│   ├── components/
│   │   ├── ui/                          # primitives
│   │   └── inventory/                   # list, filters, aging badge, dialogs
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts                # envelope unwrap
│   │   │   ├── types.ts                 # Api* snake_case
│   │   │   └── adapters.ts              # toVehicleVM, …
│   │   └── data.ts                      # TanStack Query hooks
│   └── styles.css
└── …
```
