-- Rename long_description to description and remove unused fields
-- This migration aligns the database schema with Apple App Store Connect API

-- Rename long_description to description
ALTER TABLE app_localizations RENAME COLUMN long_description TO description;

-- Remove unused fields
ALTER TABLE app_localizations DROP COLUMN IF EXISTS download_description;
ALTER TABLE app_localizations DROP COLUMN IF EXISTS short_description;