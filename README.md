# Inventory Dashboard

Dealership **inventory dashboard** assessment: filter vehicles, identify **aging stock** (in stock more than 90 days), and record a **proposed action**.

## Stack

| Layer | Choice |
|-------|--------|
| Backend | Golang (modular monolith, REST) |
| Frontend | React + TypeScript |
| Database | PostgreSQL |
| Images | MinIO (private; streamed by API) |
| Local run | Docker Compose — API and UI on port 8080 |

## Where to read

- **Words:** [GLOSSARY.md](GLOSSARY.md)
- **Decisions:** [docs/adr/](docs/adr/) (0001–0009)
- **System shape:** [docs/architecture/README.md](docs/architecture/README.md)
- **Behavior:** [.scratch/inventory-dashboard/spec.md](.scratch/inventory-dashboard/spec.md)
- **Agent entry:** [AGENTS.md](AGENTS.md)
- **Cursor rules:** `.cursor/rules/`. `architecture.mdc` and `frontend-integration.mdc` are always on. The others apply by file.

## Brief

See [initial-prompt.md](initial-prompt.md).

## Run

Docker Compose serves the API and the built screen together on port 8080. PostgreSQL is published on port 5432. MinIO stays on the Compose network and is not published to the host. The MinIO image is built from the last community release that still publishes a binary, because Docker Hub no longer serves `minio/minio`. There is no sign-in.

```bash
docker compose up --build
```

Open http://localhost:8080. `GET /health` is the liveness check, outside `/api/v1`.

During UI development the screen can run on its own and still call the API on the same paths. Leave Compose up, then from `frontend/`:

```bash
npm install
npm run dev
```

Vite proxies `/api` and `/health` to `http://localhost:8080`. See [frontend/README.md](frontend/README.md).

## Status

The inventory screen shell and `GET /health` are served on port 8080. Catalog, vehicles, and primary images are not built yet.
