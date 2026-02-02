-- Rollback: rename whats_new back to release_notes
ALTER TABLE app_localizations RENAME COLUMN whats_new TO release_notes;