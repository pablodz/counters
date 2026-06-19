CREATE TABLE IF NOT EXISTS post_metrics (
    content_id VARCHAR(255) NOT NULL,
    content_type VARCHAR(50) NOT NULL,

    views_count INTEGER DEFAULT 0,
    likes_count INTEGER DEFAULT 0,
    shares_count INTEGER DEFAULT 0,

    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (content_id, content_type)
);

CREATE INDEX IF NOT EXISTS idx_metrics_trending ON post_metrics (updated_at DESC, views_count DESC);

CREATE INDEX IF NOT EXISTS idx_metrics_type_trending ON post_metrics (content_type, updated_at DESC, views_count DESC);
