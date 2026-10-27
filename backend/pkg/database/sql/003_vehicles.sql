CREATE TABLE IF NOT EXISTS vehicle_vins (
    vin TEXT PRIMARY KEY,
    vehicle_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS vehicles (
    id BIGINT NOT NULL,
    vin TEXT NOT NULL,
    model_id INTEGER NOT NULL REFERENCES models (id),
    entered_on DATE NOT NULL,
    stocked_out_on DATE,
    proposed_action TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, created_at),
    CONSTRAINT vehicles_proposed_action_known CHECK (
        proposed_action IS NULL
        OR proposed_action IN ('price_reduction_planned', 'wholesale_planned', 'hold')
    )
) PARTITION BY RANGE (created_at);

DO $$
DECLARE
    month_start DATE := DATE '2026-01-01';
    month_stop DATE := DATE '2028-01-01';
BEGIN
    WHILE month_start < month_stop LOOP
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF vehicles FOR VALUES FROM (%L) TO (%L)',
            'vehicles_' || to_char(month_start, 'YYYY_MM'),
            month_start,
            (month_start + INTERVAL '1 month')::date
        );
        month_start := (month_start + INTERVAL '1 month')::date;
    END LOOP;
END $$;

CREATE INDEX IF NOT EXISTS vehicles_model_id_idx ON vehicles (model_id);
CREATE INDEX IF NOT EXISTS vehicles_entered_on_idx ON vehicles (entered_on);
CREATE INDEX IF NOT EXISTS vehicles_stocked_out_on_idx ON vehicles (stocked_out_on);
CREATE INDEX IF NOT EXISTS vehicles_vin_idx ON vehicles (vin);
