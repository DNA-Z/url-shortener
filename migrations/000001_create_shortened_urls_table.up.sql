-- migrations/000001_create_shortened_urls_table.up.sql

CREATE TABLE urls (
                        id uuid PRIMARY KEY,
                        short_url VARCHAR(255) NOT NULL,
                        original_url TEXT NOT NULL
);

CREATE INDEX idx_short_url ON urls(short_url);
CREATE INDEX idx_original_url ON urls(original_url);