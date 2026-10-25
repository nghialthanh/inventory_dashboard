# Database Schema

PostgreSQL is the system of record. MinIO holds primary image bytes; Postgres stores only the object key ([ADR-0003](../adr/0003-primary-image-in-minio.md)).

**Status:** planned — align migrations with this doc when scaffolding `backend/`.

## Tables

### `makes`

| Column | Type | Notes |
|--------|------|-------|
| `id` | UUID PK | Stable catalog id |
| `name` | TEXT NOT NULL | Unique display name |
| `created_at` | TIMESTAMPTZ | |

### `models`

| Column | Type | Notes |
|--------|------|-------|
| `id` | UUID PK | |
| `make_id` | UUID NOT NULL FK → `makes.id` | Model belongs to one make |
| `name` | TEXT NOT NULL | Unique per make `(make_id, name)` |
| `created_at` | TIMESTAMPTZ | |

Daily use does not edit the catalog; missing models are added via seed ([ADR-0002](../adr/0002-model-catalog.md)).

### `vehicles`

| Column | Type | Notes |
|--------|------|-------|
| `id` | UUID PK | Opaque API identity ([ADR-0006](../adr/0006-uuid-vehicle-identity.md)) |
| `vin` | TEXT NOT NULL | Unique; never reused after stock-out ([ADR-0004](../adr/0004-stocked-out-is-terminal.md)) |
| `model_id` | UUID NOT NULL FK → `models.id` | No free-text make/model |
| `entered_on` | DATE NOT NULL | Inventory age starts here |
| `stocked_out_on` | DATE NULL | NULL ⇒ in stock; set ⇒ stocked out ([ADR-0005](../adr/0005-stock-state-from-stocked-out-on.md)) |
| `proposed_action` | TEXT NULL | `price_reduction_planned` \| `wholesale_planned` \| `hold` |
| `primary_image_key` | TEXT NULL | MinIO object key; not a URL |
| `created_at` | TIMESTAMPTZ | |
| `updated_at` | TIMESTAMPTZ | |

**Do not add** a `stock_state` column — derive it from `stocked_out_on`.

Vehicles that leave inventory are retained ([ADR-0001](../adr/0001-retain-stocked-out-vehicles.md)). Prefer status/date lifecycle over hard delete.

## Derived fields (API / logic, not columns)

| Field | Definition |
|-------|------------|
| `stock_state` | `in_stock` if `stocked_out_on` IS NULL, else `stocked_out` |
| `inventory_age` | Calendar days from `entered_on` to “today” (in stock) or to `stocked_out_on` (stocked out) |
| `is_aging` | `stock_state = in_stock` AND `inventory_age > 90` |
| `has_primary_image` | `primary_image_key IS NOT NULL` |

Aging threshold constant: **90** (document in `constant` + tests).

## Indexes

| Index | Purpose |
|-------|---------|
| `UNIQUE (vin)` | VIN identity |
| `UNIQUE (makes.name)` | Catalog cleanliness |
| `UNIQUE (make_id, name)` on `models` | Model uniqueness per make |
| `(model_id)` on `vehicles` | Filter by model/make (join) |
| `(entered_on)` | Age-oriented queries if needed |
| `(stocked_out_on)` | In-stock filter (`IS NULL`) / stocked-out-on day |

## Migrations & seed

1. Migrate `makes`, `models`, `vehicles`.
2. Seed a realistic make/model catalog for demo filters.
3. Optional demo vehicles spanning ages around the 90-day boundary for aging demos.

## MinIO

- Bucket private to the compose network.
- Key scheme example: `vehicles/{vehicle_id}/primary` (exact scheme owned by `logic/media`).
- Deleting/replacing image updates `primary_image_key` in the same use case as object mutate.
