-- Revert the changes: rename description back to long_description and add back the fields

-- Rename description back to long_description
ALTER TABLE app_localizations RENAME COLUMN description TO long_description;

-- Add back the removed fields
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS download_description TEXT;
ALTER TABLE app_localizations ADD COLUMN IF NOT EXISTS short_description VARCHAR(80);