-- Remove version_state field from app_localizations table
ALTER TABLE app_localizations DROP COLUMN IF EXISTS version_state;