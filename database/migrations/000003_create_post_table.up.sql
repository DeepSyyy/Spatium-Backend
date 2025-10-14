CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Migration Up: create posts table
CREATE TABLE IF NOT EXISTS posts (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    user_internal_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    ai_response TEXT,
    mood_tag_internal_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_posts_user
        FOREIGN KEY (user_internal_id)
        REFERENCES users (internal_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_posts_mood
        FOREIGN KEY (mood_tag_internal_id)
        REFERENCES mood_tags (internal_id)
        ON DELETE SET NULL
);
