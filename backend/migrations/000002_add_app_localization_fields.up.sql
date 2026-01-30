-- Add promotional_text and what_to_test fields to app_localizations table
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS promotional_text TEXT;
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS what_to_test TEXT;
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS locale VARCHAR(50);