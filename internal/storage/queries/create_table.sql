CREATE TABLE IF NOT EXISTS urls (
    uuid UUID PRIMARY KEY,
    user_id uuid NULL,
    short_url VARCHAR(255) NULL,
    original_url TEXT NULL,
    id_deleted BOOLEAN DEFAULT FALSE
);