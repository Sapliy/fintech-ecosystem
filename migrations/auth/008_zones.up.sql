CREATE TABLE IF NOT EXISTS zones (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL,
    name TEXT NOT NULL,
    mode TEXT NOT NULL CHECK (mode IN ('test', 'live')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add indexes for common lookups
CREATE INDEX idx_zones_org_id ON zones(org_id);

-- Alter api_keys to belong to a zone
ALTER TABLE api_keys ADD COLUMN zone_id TEXT REFERENCES zones(id);
ALTER TABLE api_keys ADD COLUMN mode TEXT CHECK (mode IN ('test', 'live'));

-- Create index for API key lookups by zone
CREATE INDEX idx_api_keys_zone_id ON api_keys(zone_id);
