-- migrations/000003_add_is_deleted_column.up.sql
ALTER TABLE urls ADD COLUMN is_deleted boolean DEFAULT false NOT NULL;