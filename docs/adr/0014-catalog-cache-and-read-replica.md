# Catalog lists use Redis cache-aside; catalog SQL reads use a Postgres replica

Make and model list responses are hot, tiny, and unchanged by daily API use ([ADR-0002](0002-model-catalog.md)). The API caches those list payloads in Redis (cache-aside, ~1h TTL) and fails open to PostgreSQL when Redis errors. Catalog `ListMakes` / `ListModels` queries use a read DSN (streaming replica in Compose); vehicle writes and `ModelMakeID` stay on the primary so record/correct keep read-your-writes. Process shutdown drains HTTP for a configured timeout, then closes pools. Prometheus metrics on `GET /metrics` expose cache hit/miss/error and HTTP RED-style signals.

## Considered options

- In-process memory only. Enough for one API replica; rejected for the shared-cache and multi-instance story.
- Cache vehicle lists. Rejected: inventory changes often; invalidation would dominate the demo.
- Require Redis (fail closed). Rejected: listing makes/models must survive cache outage.
