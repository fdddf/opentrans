-- Rename release_notes column to whats_new to match Apple API field name
ALTER TABLE app_localizations RENAME COLUMN release_notes TO whats_new;