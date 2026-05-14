-- migrations/000002_add_user_id_column.up.sql
ALTER TABLE urls ADD COLUMN user_id uuid NULL;
CREATE INDEX idx_user_id ON urls(uuid);