-- Drop indexes
DROP INDEX IF EXISTS idx_blocks_blocker_id;
DROP INDEX IF EXISTS idx_blocks_blocked_user_id;

-- Drop blocks table
DROP TABLE IF EXISTS blocks;
