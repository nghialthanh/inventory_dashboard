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
| `VALIDATION_ERROR` | 400 | Bad dates, VIN, enums, unknown model, illegal filter pair. Also a string, fraction, zero, negative, or non-canonical address, a make or model above 2147483647, or a vehicle address above 9007199254740991 |
| `NOT_FOUND` | 404 | A positive in-range vehicle address that matches no vehicle |
| `CONFLICT` | 409 | Duplicate VIN; correction after stock-out; clearing `stocked_out_on`; proposed action when the vehicle is not aging stock |
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

Path id is the vehicle address, a positive integer, not the VIN ([ADR-0011](../adr/0011-integer-addresses.md)). Make and model addresses are 32-bit integers. A vehicle address above 9007199254740991 is a validation error.

List and detail share one shape, with embedded make and model:

```json
{
  "id": 5,
  "vin": "1HGCM82633A004352",
  "make": { "id": 2, "name": "Honda" },
  "model": { "id": 4, "name": "Civic" },
  "entered_on": "2026-01-10",
  "stocked_out_on": null,
  "stock_state": "in_stock",
  "inventory_age": 95,
  "is_aging": true,
  "proposed_action": "price_reduction_planned"
}
```

`stock_state` is `in_stock` or `stocked_out`. `is_aging` is in stock and `inventory_age` > 90.

### List

`GET /api/v1/vehicles` — **implemented** for make, model, stock state, stocked-out on, aging, and paging.

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
{ "vin": "1HGCM82633A004352", "model_id": 4, "entered_on": "2026-01-10" }
```

The vehicle starts in stock, with no proposed action. `entered_on` is today or earlier. `model_id` must exist.

**VIN:** trim, store uppercase, 17 characters, digits and letters except I, O, and Q. No check digit. Unique, including against a stocked-out vehicle.

### Correct

`PATCH /api/v1/vehicles/:id` — **implemented**. It accepts a non-empty subset of `vin`, `model_id`, and `entered_on` while the vehicle is in stock. The vehicle address stays. The same command on a stocked-out vehicle is a conflict: the correction is rejected because the vehicle is stocked out. This route rejects `proposed_action` and `stocked_out_on`.

### Commands

| Method | Path | Body | Rule |
|--------|------|------|------|
| POST | `/api/v1/vehicles/:id/stock-out` | `{ "stocked_out_on": "YYYY-MM-DD" }` | **implemented**. First call stocks out. A later call replaces the date. On or after `entered_on`, on or before today. Empty is rejected. Clearing it is a conflict: the departure stands. This is the only way to change `stocked_out_on` |
| PUT | `/api/v1/vehicles/:id/proposed-action` | `{ "proposed_action": "wholesale_planned" }` | **implemented**. `price_reduction_planned`, `wholesale_planned`, `hold`, or empty to clear. Accepted only while the vehicle is aging stock. After stock-out the stored value remains on read and writes conflict |

## Static UI

The same process serves the SPA. Unknown non-API paths fall back to `index.html`.
