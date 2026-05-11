-- migrations/000002_add_user_id_column.down

DROP INDEX IF EXISTS idx_user_id;

ALTER TABLE urls DROP COLUMN IF EXISTS user_id;