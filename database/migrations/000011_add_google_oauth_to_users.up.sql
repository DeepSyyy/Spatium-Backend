-- Add Google OAuth fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS photo_url TEXT;

-- Make recovery_code nullable since Google users don't need it
ALTER TABLE users ALTER COLUMN recovery_code DROP NOT NULL;
