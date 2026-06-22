CREATE TABLE IF NOT EXISTS item_interactions_hourly (
    item_id TEXT NOT NULL,
    item_type TEXT NOT NULL,
    event_type TEXT NOT NULL,
    period_hour_unix INTEGER NOT NULL,
    total_count INTEGER DEFAULT 1,
    PRIMARY KEY (item_id, item_type, event_type, period_hour_unix)
);

CREATE TABLE IF NOT EXISTS user_item_interactions_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    user_type TEXT NOT NULL,
    item_id TEXT NOT NULL,
    item_type TEXT NOT NULL,
    event_type TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_item_log ON user_item_interactions_log (user_id, item_id);
