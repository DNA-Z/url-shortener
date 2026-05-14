-- migrations/000001_create_shortened_urls_table.down.sql
DROP INDEX IF EXISTS idx_original_url;
DROP INDEX IF EXISTS idx_short_url;
DROP INDEX IF EXISTS idx_user_id;
DROP TABLE IF EXISTS urls;