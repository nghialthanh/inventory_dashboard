CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY,
    vin TEXT NOT NULL UNIQUE,
    model_id UUID NOT NULL REFERENCES models (id),
    entered_on DATE NOT NULL,
    stocked_out_on DATE,
    proposed_action TEXT,
    primary_image_key TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT vehicles_proposed_action_known CHECK (
        proposed_action IS NULL
        OR proposed_action IN ('price_reduction_planned', 'wholesale_planned', 'hold')
    )
);

CREATE INDEX IF NOT EXISTS vehicles_model_id_idx ON vehicles (model_id);
CREATE INDEX IF NOT EXISTS vehicles_entered_on_idx ON vehicles (entered_on);
CREATE INDEX IF NOT EXISTS vehicles_stocked_out_on_idx ON vehicles (stocked_out_on);
