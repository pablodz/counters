CREATE TABLE IF NOT EXISTS item_snapshots (
    item_id VARCHAR(255) NOT NULL,
    item_type VARCHAR(50) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    period_hour_unix INTEGER NOT NULL,
    total_count INTEGER DEFAULT 0,
    PRIMARY KEY (item_id, item_type, event_type, period_hour_unix)
);

CREATE INDEX IF NOT EXISTS idx_snapshots_query
    ON item_snapshots (item_id, item_type, event_type, period_hour_unix);
