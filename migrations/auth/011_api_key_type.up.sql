-- Add type column to api_keys
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS type VARCHAR(20) DEFAULT 'secret'; -- secret or publishable

-- Update some existing keys if necessary
-- UPDATE api_keys SET type = 'secret' WHERE key_prefix LIKE 'sk_%';
