-- Add sync metadata fields to app_localizations table
-- These fields were defined in the original schema but may be missing in some databases
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS synced_at TIMESTAMP WITHOUT TIME ZONE;
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS source VARCHAR(20) NOT NULL DEFAULT 'local';
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS sync_status VARCHAR(20) NOT NULL DEFAULT 'pending';