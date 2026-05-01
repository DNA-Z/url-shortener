-- migrations/000001_create_shortened_urls_table.up.sql
-- Создание таблицы сокращенных URL
CREATE TABLE shortened_urls (
                        id uuid PRIMARY KEY,
                        short_url VARCHAR(255) NOT NULL,
                        original_url TEXT NOT NULL
);

-- Базовый индекс для поиска по сокращенному URL
CREATE INDEX idx_short_url ON shortened_urls(short_url);

-- Индекс для поиска по оригинальному URL
CREATE INDEX idx_original_url ON shortened_urls(original_url);