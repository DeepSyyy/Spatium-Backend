-- Create blocks table for user blocks
CREATE TABLE IF NOT EXISTS blocks (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    blocker_id BIGINT NOT NULL REFERENCES users(internal_id) ON DELETE CASCADE,
    blocked_user_id BIGINT NOT NULL REFERENCES users(internal_id) ON DELETE CASCADE,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(blocker_id, blocked_user_id)
);

-- Create indexes for efficient querying
CREATE INDEX idx_blocks_blocker_id ON blocks(blocker_id);
CREATE INDEX idx_blocks_blocked_user_id ON blocks(blocked_user_id);
