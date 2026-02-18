-- Remove added columns from payment_intents
ALTER TABLE payment_intents DROP COLUMN IF EXISTS description;
ALTER TABLE payment_intents DROP COLUMN IF EXISTS zone_id;
ALTER TABLE payment_intents DROP COLUMN IF EXISTS mode;
