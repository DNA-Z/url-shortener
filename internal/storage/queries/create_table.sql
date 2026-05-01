CREATE TABLE IF NOT EXISTS urls (
    uuid UUID PRIMARY KEY,
    short_url TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL
);