# Database Schema

PostgreSQL is the system of record ([ADR-0010](../adr/0010-no-object-storage.md)).

**Status:** `makes`, `models`, `vehicle_vins`, and partitioned `vehicles` are migrated. Makes, models, and demo vehicles (ages 0, 90, 91, plus one stocked-out stay longer than 90 days) are seeded.

## Tables

### `makes`

| Column | Type | Notes |
|--------|------|-------|
| `id` | INTEGER identity PK | Catalog address. Seed pins Ford 1, Honda 2, Toyota 3 ([ADR-0011](../adr/0011-integer-addresses.md)) |
| `name` | TEXT NOT NULL | Unique display name |
| `created_at` | TIMESTAMPTZ | |

### `models`

| Column | Type | Notes |
|--------|------|-------|
| `id` | INTEGER identity PK | Seed pins F-150 1 through Corolla 6 |
| `make_id` | INTEGER NOT NULL FK → `makes.id` | Model belongs to one make |
| `name` | TEXT NOT NULL | Unique per make `(make_id, name)` |
| `created_at` | TIMESTAMPTZ | |

Daily use does not edit the catalog; missing models are added via seed ([ADR-0002](../adr/0002-model-catalog.md)).

### `vehicle_vins`

| Column | Type | Notes |
|--------|------|-------|
| `vin` | TEXT PK | Global VIN uniqueness across months ([ADR-0004](../adr/0004-stocked-out-is-terminal.md), [ADR-0013](../adr/0013-snowflake-vehicle-addresses.md)) |
| `vehicle_id` | BIGINT NOT NULL | Vehicle address |

### `vehicles`

RANGE-partitioned by month on `created_at` (partitions `vehicles_2026_01` … `vehicles_2027_12`). No DEFAULT partition ([ADR-0013](../adr/0013-snowflake-vehicle-addresses.md)).

| Column | Type | Notes |
|--------|------|-------|
| `id` | BIGINT NOT NULL | Vehicle address (app-minted snowflake). JSON string on the wire |
| `vin` | TEXT NOT NULL | Also registered in `vehicle_vins`. Shape is in [api-outline.md](api-outline.md) |
| `model_id` | INTEGER NOT NULL FK → `models.id` | No free-text make/model |
| `entered_on` | DATE NOT NULL | Inventory age starts here (correctable; not the partition key) |
| `stocked_out_on` | DATE NULL | NULL ⇒ in stock; set ⇒ stocked out ([ADR-0005](../adr/0005-stock-state-from-stocked-out-on.md)) |
| `proposed_action` | TEXT NULL | `price_reduction_planned` \| `wholesale_planned` \| `hold` |
| `created_at` | TIMESTAMPTZ NOT NULL | Mint instant; partition key; matches snowflake time |
| `updated_at` | TIMESTAMPTZ | |

Primary key: `(id, created_at)`.

Derive stock state from `stocked_out_on`. Keep the row after stock-out ([ADR-0001](../adr/0001-retain-stocked-out-vehicles.md)). There is no delete-vehicle operation.

## Derived fields (API / logic, not columns)

| Field | Definition |
|-------|------------|
| `stock_state` | `in_stock` if `stocked_out_on` IS NULL, else `stocked_out` |
| `inventory_age` | Calendar days from `entered_on` to today while in stock, or to `stocked_out_on` when stocked out. The entered-on day is 0 |
| `is_aging` | `stock_state = in_stock` AND `inventory_age > 90` |

Aging threshold constant: **90** (document in `constant` + tests).

## Indexes

| Index | Purpose |
|-------|---------|
| `PRIMARY KEY (vin)` on `vehicle_vins` | Global VIN identity |
| `UNIQUE (makes.name)` | Catalog cleanliness |
| `UNIQUE (make_id, name)` on `models` | Model uniqueness per make |
| `(model_id)` on `vehicles` | Filter by model/make (join) |
| `(entered_on)` | Age-oriented queries if needed |
| `(stocked_out_on)` | In-stock filter (`IS NULL`) / stocked-out-on day |
| `(vin)` on `vehicles` | Lookup aid (not the global unique) |

## Migrations & seed

1. Migrate `makes`, `models`, `vehicle_vins`, partitioned `vehicles`.
2. Seed a realistic make/model catalog for demo filters.
3. Seed demo vehicles from the UTC date at seed time: inventory age 0, 90, and 91, plus one stocked-out vehicle whose stay was longer than 90 days. Those four vehicles use fixed snowflake addresses for worker 0 (see `004_seed_vehicles.sql`). Drop the Compose volume before applying this schema over an older layout.
