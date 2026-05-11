-- migrations/000003_add_is_deleted_column.down

ALTER TABLE urls DROP COLUMN IF EXISTS is_deleted;