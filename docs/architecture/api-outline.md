# API outline

Base: `/api/v1`. The UI and the API share one origin ([ADR-0009](../adr/0009-one-origin-under-api.md)).

**Status:** planned. Mark a row **implemented** when the route exists.

No authentication. Dates are calendar dates. "Today" is the UTC date of a clock injected into the API; tests set that clock.

## Envelope

```json
{ "code": 200, "message": "ok", "data": {}, "meta": { "page": 1, "page_size": 20, "total": 100 } }
```

Errors put `error_code` and optional `fields` in `data`. JSON is snake_case.

| `error_code` | HTTP | When |
|--------------|------|------|
| `VALIDATION_ERROR` | 400 | Bad dates, VIN, enums, unknown model, illegal filter pair |
| `NOT_FOUND` | 404 | Unknown vehicle, or a primary image that is not stored |
| `CONFLICT` | 409 | Duplicate VIN; correction or image change after stock-out; clearing `stocked_out_on`; proposed action when the vehicle is not aging stock |
| `INTERNAL_ERROR` | 500 | Unexpected |

## Health

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness. Outside `/api/v1`. **implemented** |

## Catalog

Read-only. Seed owns writes ([ADR-0002](../adr/0002-model-catalog.md)).

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/catalog/makes` | `{ id, name }`. **implemented** |
| GET | `/api/v1/catalog/models?make_id=` | Models `{ id, make_id, name }`. `make_id` is optional. **implemented** |

## Vehicles

Path id is the vehicle UUID, not the VIN ([ADR-0006](../adr/0006-uuid-vehicle-identity.md)).

List and detail share one shape, with embedded make and model:

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

`stock_state` is `in_stock` or `stocked_out`. `is_aging` is in stock and `inventory_age` > 90. No MinIO URL or object key.

### List

`GET /api/v1/vehicles` — **implemented** for make, model, aging, and paging. `stock_state` and `stocked_out_on` filters are not implemented yet.

| Param | Meaning |
|-------|---------|
| `make_id` | Filter through the vehicle's model |
| `model_id` | Filter by model |
| `stock_state` | `in_stock` or `stocked_out` |
| `aging` | `true` only. Any other value is a validation error. Omit to leave the filter off |
| `stocked_out_on` | Exact calendar day |
| `page` | Default 1 |
| `page_size` | Default 20, maximum 100 |

Sort is `entered_on` ascending, then `id` ascending. A page past the end is an empty list with the true `total`. If both `make_id` and `model_id` are set and that model is not of that make, respond `VALIDATION_ERROR`.

### Read

`GET /api/v1/vehicles/:id` — **implemented**

### Create

`POST /api/v1/vehicles` — **implemented**

```json
{ "vin": "1HGCM82633A004352", "model_id": "…", "entered_on": "2026-01-10" }
```

The vehicle starts in stock, with no proposed action and no primary image. `entered_on` is today or earlier. `model_id` must exist.

**VIN:** trim, store uppercase, 17 characters, digits and letters except I, O, and Q. No check digit. Unique, including against a stocked-out vehicle.

### Correct

`PATCH /api/v1/vehicles/:id` accepts a non-empty subset of `vin`, `model_id`, and `entered_on` while the vehicle is in stock. The UUID stays. The same command on a stocked-out vehicle is a conflict. This route rejects `proposed_action`, image, and `stocked_out_on`.

### Commands

| Method | Path | Body | Rule |
|--------|------|------|------|
| POST | `/api/v1/vehicles/:id/stock-out` | `{ "stocked_out_on": "YYYY-MM-DD" }` | First call stocks out. A later call replaces the date. On or after `entered_on`, on or before today. Empty is rejected. This is the only way to change `stocked_out_on` |
| PUT | `/api/v1/vehicles/:id/proposed-action` | `{ "proposed_action": "wholesale_planned" }` | `price_reduction_planned`, `wholesale_planned`, `hold`, or empty to clear. Accepted only while the vehicle is aging stock. After stock-out the stored value remains on read and writes conflict |
| PUT | `/api/v1/vehicles/:id/primary-image` | multipart field `file` | JPEG, PNG, or WebP, at most 5 MiB. In stock only. Replaces the previous object |
| DELETE | `/api/v1/vehicles/:id/primary-image` | — | Removes the object and clears the key. In stock only |
| GET | `/api/v1/vehicles/:id/primary-image` | — | Streams the bytes ([ADR-0008](../adr/0008-api-streams-images.md)). Not found when there is no key or no object. Still streams after stock-out |

## Static UI

The same process serves the SPA. Unknown non-API paths fall back to `index.html`.
