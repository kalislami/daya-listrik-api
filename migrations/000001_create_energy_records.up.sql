CREATE TABLE IF NOT EXISTS energy_records (
    id SERIAL PRIMARY KEY,
    date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usage DOUBLE PRECISION NOT NULL,
    device VARCHAR(100) NOT NULL,
    duration REAL NOT NULL DEFAULT 1
);

-- Existing installations used a nullable duration column.
ALTER TABLE energy_records ADD COLUMN IF NOT EXISTS duration REAL DEFAULT 1;
UPDATE energy_records SET duration = 1 WHERE duration IS NULL;
ALTER TABLE energy_records ALTER COLUMN duration SET DEFAULT 1;
ALTER TABLE energy_records ALTER COLUMN duration SET NOT NULL;
