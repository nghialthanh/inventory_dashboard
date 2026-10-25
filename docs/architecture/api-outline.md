# API Outline

Base URL: `/api/v1`  
UI + API share one origin in Compose ([ADR-0009](../adr/0009-one-origin-under-api.md)).

**Status:** planned contract — mark endpoints **implemented** in this file when routes land.

Authentication: assessment-simple (none or a single shared credential). Do not invent multi-tenant RBAC unless the brief expands.

## Response envelope

```json
{
  "code": 200,
  "message": "ok",
  "data": {},
  "meta": { "page": 1, "page_size": 20, "total": 100 }
}
```

Error `data` may contain `error_code` and field errors:

```json
{
  "code": 409,
  "message": "vin already exists",
  "data": {
    "error_code": "CONFLICT",
    "fields": { "vin": "duplicate" }
  }
}
```

JSON fields: **snake_case**.

---

## Health

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness (outside `/api/v1` or under it — pick one and stick to it) |

---

## Catalog

Read-only for daily use; seed owns writes ([ADR-0002](../adr/0002-model-catalog.md)).

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/catalog/makes` | List makes `{ id, name }` |
| GET | `/api/v1/catalog/models?make_id=` | List models; optional filter by make |

---

## Vehicles

Path id is the vehicle **UUID**, not the VIN ([ADR-0006](../adr/0006-uuid-vehicle-identity.md)).

### List & read

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/vehicles` | Filterable list |
| GET | `/api/v1/vehicles/:id` | Detail |

**Query params (list):**

| Param | Meaning |
|-------|---------|
| `make_id` | Filter via model → make |
| `model_id` | Filter by model |
| `stock_state` | `in_stock` \| `stocked_out` (optional) |
| `aging` | `true` → only aging stock (in stock and age > 90) |
| `stocked_out_on` | Exact calendar day (managers reviewing departures) |
| `page`, `page_size` | Pagination (`page_size` capped) |

**List/detail item (illustrative):**

```json
{
  "id": "018f…",
  "vin": "1HGCM82633A004352",
  "make": { "id": "…", "name": "Honda" },
  "model": { "id": "…", "name": "Civic" },
  "entered_on": "2026-01-10",
  "stocked_out_on": null,
  "stock_state": "in_stock",
  "inventory_age": 95,
  "is_aging": true,
  "proposed_action": "price_reduction_planned",
  "has_primary_image": true
}
```

### Create & correct

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/vehicles` | Create in-stock vehicle |
| PATCH | `/api/v1/vehicles/:id` | Correct fields allowed for current stock state |

**POST body (illustrative):**

```json
{
  "vin": "1HGCM82633A004352",
  "model_id": "…",
  "entered_on": "2026-01-10"
}
```

**PATCH rules (summary):**

- While **in stock**: VIN correction, and other allowed corrections (document exact allow-list in code + here when implemented).
- While **stocked out**: do **not** clear `stocked_out_on` or reopen; stocked-out-on date correction belongs on the stock-out command or a dedicated correction — prefer one place and document it ([ADR-0004](../adr/0004-stocked-out-is-terminal.md)).
- Do **not** accept `proposed_action`, image bytes, or stock-out fields on this PATCH ([ADR-0007](../adr/0007-separate-commands.md)).

### Commands (separate routes)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/vehicles/:id/stock-out` | Set `stocked_out_on`; terminal transition |
| PUT | `/api/v1/vehicles/:id/proposed-action` | Set/replace proposed action (`null` to clear if product allows) |
| PUT | `/api/v1/vehicles/:id/primary-image` | Upload/replace image (`multipart` or raw — pick one) |
| DELETE | `/api/v1/vehicles/:id/primary-image` | Remove image + clear key |
| GET | `/api/v1/vehicles/:id/primary-image` | Stream image bytes ([ADR-0008](../adr/0008-api-streams-images.md)) |

**Stock-out body:**

```json
{ "stocked_out_on": "2026-10-25" }
```

**Proposed action body:**

```json
{ "proposed_action": "wholesale_planned" }
```

Allowed values: `price_reduction_planned`, `wholesale_planned`, `hold`.

Vehicle JSON never includes MinIO URLs or object keys — only `has_primary_image`.

---

## Error codes (machine-readable)

| `error_code` | Typical HTTP | When |
|--------------|--------------|------|
| `VALIDATION_ERROR` | 400 | Bad input |
| `NOT_FOUND` | 404 | Unknown id; missing image when GET requires one |
| `CONFLICT` | 409 | Duplicate VIN; illegal stock transition |
| `INTERNAL_ERROR` | 500 | Unexpected |

---

## Static UI

Non-API routes served by the same process (SPA fallback to `index.html`) so the browser stays on one origin.
