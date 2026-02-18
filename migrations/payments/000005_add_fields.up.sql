-- Add missing columns to payment_intents
ALTER TABLE payment_intents ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE payment_intents ADD COLUMN IF NOT EXISTS zone_id VARCHAR(50);
ALTER TABLE payment_intents ADD COLUMN IF NOT EXISTS mode VARCHAR(10);
