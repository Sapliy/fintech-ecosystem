-- Increase status column length
ALTER TABLE payment_intents ALTER COLUMN status TYPE VARCHAR(50);
