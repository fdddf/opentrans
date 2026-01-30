-- Remove promotional_text and what_to_test fields from app_localizations table
ALTER TABLE app_localizations DROP COLUMN IF EXISTS promotional_text;
ALTER TABLE app_localizations DROP COLUMN IF EXISTS what_to_test;
ALTER TABLE app_localizations DROP COLUMN IF EXISTS locale;