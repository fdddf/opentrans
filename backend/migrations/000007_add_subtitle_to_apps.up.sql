-- Add subtitle field to apps table
ALTER TABLE apps ADD COLUMN IF NOT EXISTS subtitle TEXT;