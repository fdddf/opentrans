-- Remove sync metadata fields from app_localizations table
ALTER TABLE app_localizations DROP COLUMN IF EXISTS synced_at;
ALTER TABLE app_localizations DROP COLUMN IF EXISTS source;
ALTER TABLE app_localizations DROP COLUMN IF EXISTS sync_status;