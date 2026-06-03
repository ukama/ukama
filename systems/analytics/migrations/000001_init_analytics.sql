CREATE TABLE IF NOT EXISTS analytics_event_log (
    id BIGSERIAL PRIMARY KEY,
    routing_key TEXT NOT NULL,
    event_hash TEXT NOT NULL,
    org_name TEXT,
    received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ,
    UNIQUE(routing_key, event_hash)
);

CREATE TABLE IF NOT EXISTS analytics_event_error (
    id BIGSERIAL PRIMARY KEY,
    routing_key TEXT NOT NULL,
    event_hash TEXT NOT NULL,
    error TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS analytics_refresh_state (
    source TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    last_started_at TIMESTAMPTZ,
    last_completed_at TIMESTAMPTZ,
    error TEXT
);
