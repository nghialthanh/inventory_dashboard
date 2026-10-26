-- Bulk demo data: 1100 in-stock / mixed vehicles, of which 150 are aging stock
-- (in stock, inventory age > 90). Safe to re-run: VINs use a LOAD prefix and
-- ON CONFLICT DO NOTHING.
--
-- Usage (Compose Postgres published on 5432):
--   docker compose exec -T postgres psql -U inventory -d inventory < scripts/seed_bulk_vehicles.sql
-- Or from the host:
--   PGPASSWORD=inventory psql -h localhost -U inventory -d inventory -f scripts/seed_bulk_vehicles.sql

WITH model_ids AS (
    SELECT ARRAY[
        1, 2, 3, 4, 5, 6
    ] AS ids
),
today AS (
    SELECT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date AS d
),
-- VIN alphabet excludes I, O, Q (matches API NormalizeVIN).
alphabet AS (
    SELECT '0123456789ABCDEFGHJKLMNPRSTUVWXYZ'::text AS chars
),
rows AS (
    SELECT
        g AS n,
        -- 17-char VIN: LOAD + base-33 encoding of n (no I/O/Q)
        (
            'LOAD'
            || substr(a.chars, ((g / 33^3)::int % 33) + 1, 1)
            || substr(a.chars, ((g / 33^2)::int % 33) + 1, 1)
            || substr(a.chars, ((g / 33)::int % 33) + 1, 1)
            || substr(a.chars, (g % 33) + 1, 1)
            || lpad((g % 1000000000)::text, 9, '0')
        ) AS vin,
        m.ids[1 + ((g - 1) % array_length(m.ids, 1))] AS model_id,
        CASE
            -- 150 aging stock: in stock, entered 91..240 days ago
            WHEN g <= 150 THEN t.d - (90 + (g % 150) + 1)
            -- 50 at exactly 90 days (in stock, not aging)
            WHEN g <= 200 THEN t.d - 90
            -- 50 stocked out after a long stay (not aging)
            WHEN g <= 250 THEN t.d - 150
            -- remainder: fresh / recent in stock (age 0..89)
            ELSE t.d - ((g - 251) % 90)
        END AS entered_on,
        CASE
            WHEN g > 200 AND g <= 250 THEN t.d - 20
            ELSE NULL
        END AS stocked_out_on
    FROM generate_series(1, 1100) AS g
    CROSS JOIN model_ids m
    CROSS JOIN today t
    CROSS JOIN alphabet a
)
INSERT INTO vehicles (vin, model_id, entered_on, stocked_out_on)
SELECT vin, model_id, entered_on, stocked_out_on
FROM rows
ON CONFLICT (vin) DO NOTHING;

-- Summary
SELECT
    count(*) AS total_vehicles,
    count(*) FILTER (WHERE stocked_out_on IS NULL) AS in_stock,
    count(*) FILTER (
        WHERE stocked_out_on IS NULL
          AND ((CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - entered_on) > 90
    ) AS aging_stock
FROM vehicles;
