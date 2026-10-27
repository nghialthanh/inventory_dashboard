-- Bulk demo data: ≥510 vehicles in every monthly partition (2026-01 … 2027-12).
-- Snowflake: epoch 2024-01-01 UTC, worker 0, sequence = row-in-month - 1;
-- created_at is mid-month + 1ms per row so ids prune to the matching partition.
-- Idempotent via ON CONFLICT. Safe alongside 004_seed_vehicles.sql (those use 2026-10-01 noon).
--
--   docker compose exec -T postgres psql -U inventory -d inventory < scripts/seed_bulk_vehicles.sql

WITH params AS (
    SELECT
        510::int AS per_month,
        DATE '2026-01-01' AS first_month,
        DATE '2027-12-01' AS last_month
),
model_ids AS (
    SELECT ARRAY[1, 2, 3, 4, 5, 6] AS ids
),
today AS (
    SELECT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date AS d
),
alphabet AS (
    SELECT '0123456789ABCDEFGHJKLMNPRSTUVWXYZ'::text AS chars
),
months AS (
    SELECT
        gs::date AS month_start,
        row_number() OVER (ORDER BY gs)::int AS month_n
    FROM params p,
    LATERAL generate_series(p.first_month, p.last_month, INTERVAL '1 month') AS gs
),
rows AS (
    SELECT
        m.month_start,
        m.month_n,
        g AS n_in_month,
        ((m.month_n - 1) * p.per_month + g) AS n,
        (
            'LOAD'
            || substr(a.chars, ((((m.month_n - 1) * p.per_month + g) / (33 * 33 * 33))::int % 33) + 1, 1)
            || substr(a.chars, ((((m.month_n - 1) * p.per_month + g) / (33 * 33))::int % 33) + 1, 1)
            || substr(a.chars, ((((m.month_n - 1) * p.per_month + g) / 33)::int % 33) + 1, 1)
            || substr(a.chars, (((m.month_n - 1) * p.per_month + g) % 33) + 1, 1)
            || lpad((((m.month_n - 1) * p.per_month + g) % 1000000000)::text, 9, '0')
        ) AS vin,
        mi.ids[1 + ((((m.month_n - 1) * p.per_month + g) - 1) % array_length(mi.ids, 1))] AS model_id,
        CASE
            WHEN g <= 80 THEN t.d - (90 + (g % 80) + 1)
            WHEN g <= 100 THEN t.d - 90
            WHEN g <= 130 THEN t.d - 150
            ELSE t.d - ((g - 131) % 90)
        END AS entered_on,
        CASE
            WHEN g > 100 AND g <= 130 THEN t.d - 20
            ELSE NULL
        END AS stocked_out_on,
        (
            (m.month_start + 14)::timestamp AT TIME ZONE 'UTC'
            + ((g - 1) * INTERVAL '1 millisecond')
        ) AS created_at
    FROM months m
    CROSS JOIN params p
    CROSS JOIN generate_series(1, (SELECT per_month FROM params)) AS g
    CROSS JOIN model_ids mi
    CROSS JOIN today t
    CROSS JOIN alphabet a
),
minted AS (
    SELECT
        r.*,
        (
            ((EXTRACT(EPOCH FROM r.created_at) * 1000)::bigint - 1704067200000) << 22
        ) | ((r.n_in_month - 1) & 4095) AS vehicle_id
    FROM rows r
),
vins AS (
    INSERT INTO vehicle_vins (vin, vehicle_id)
    SELECT vin, vehicle_id FROM minted
    ON CONFLICT DO NOTHING
    RETURNING vin
)
INSERT INTO vehicles (id, vin, model_id, entered_on, stocked_out_on, created_at, updated_at)
SELECT m.vehicle_id, m.vin, m.model_id, m.entered_on, m.stocked_out_on, m.created_at, m.created_at
FROM minted m
JOIN vins USING (vin)
ON CONFLICT DO NOTHING;

SELECT
    tableoid::regclass AS partition,
    count(*) AS vehicles
FROM vehicles
GROUP BY tableoid
ORDER BY partition;

SELECT
    v.total_vehicles,
    v.in_stock,
    v.aging_stock,
    p.partition_count,
    p.min_per_partition,
    p.max_per_partition
FROM (
    SELECT
        count(*) AS total_vehicles,
        count(*) FILTER (WHERE stocked_out_on IS NULL) AS in_stock,
        count(*) FILTER (
            WHERE stocked_out_on IS NULL
              AND ((CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - entered_on) > 90
        ) AS aging_stock
    FROM vehicles
) v
CROSS JOIN (
    SELECT
        count(*)::int AS partition_count,
        min(cnt)::int AS min_per_partition,
        max(cnt)::int AS max_per_partition
    FROM (
        SELECT count(*) AS cnt
        FROM vehicles
        GROUP BY tableoid
    ) per
) p;
