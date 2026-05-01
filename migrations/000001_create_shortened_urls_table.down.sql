-- migrations/000001_create_shortened_urls_table.down.sql
-- Откат создания таблицы сокращенных URL
DROP INDEX IF EXISTS idx_original_url;
DROP INDEX IF EXISTS idx_short_url;
DROP TABLE IF EXISTS shortened_urls;