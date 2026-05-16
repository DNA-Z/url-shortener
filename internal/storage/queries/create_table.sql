CREATE TABLE IF NOT EXISTS urls (
    uuid UUID PRIMARY KEY,
    user_id UUID NULL,
    short_url VARCHAR(255) NULL,
    original_url TEXT NULL
    is_deleted BOOLEAN DEFAULT FALSE
);