-- Remove type column from api_keys
ALTER TABLE api_keys DROP COLUMN IF NOT EXISTS type;
