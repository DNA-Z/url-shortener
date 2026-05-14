CREATE TABLE IF NOT EXISTS urls (
    uuid UUID PRIMARY KEY,
    short_url VARCHAR(255) NULL,
    original_url TEXT NULL
);