-- migrations/000001_create_shortened_urls_table.up.sql
CREATE TABLE urls (
                        id uuid PRIMARY KEY,
                        user_id uuid NULL,
                        short_url VARCHAR(255) NULL,
                        original_url TEXT NULL,
                        id_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_user_id ON urls(user_id);
CREATE INDEX idx_short_url ON urls(short_url);
CREATE INDEX idx_original_url ON urls(original_url);