-- Add version field to app_localizations table
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS version VARCHAR(50);