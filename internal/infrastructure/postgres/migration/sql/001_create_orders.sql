CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    full_name TEXT NOT NULL DEFAULT '',
    followers_count INTEGER NOT NULL DEFAULT 0 CHECK (followers_count >= 0),
    status INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_username ON orders (username);
