-- Add version_state field to app_localizations table
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS version_state VARCHAR(50);