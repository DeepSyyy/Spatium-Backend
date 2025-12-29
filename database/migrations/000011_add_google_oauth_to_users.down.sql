-- Revert Google OAuth fields from users table
ALTER TABLE users DROP COLUMN IF EXISTS google_id;
ALTER TABLE users DROP COLUMN IF EXISTS email;
ALTER TABLE users DROP COLUMN IF EXISTS photo_url;

-- Make recovery_code NOT NULL again
ALTER TABLE users ALTER COLUMN recovery_code SET NOT NULL;
