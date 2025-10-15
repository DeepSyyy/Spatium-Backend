CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS comments (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID DEFAULT gen_random_uuid() UNIQUE,
    post_internal_id BIGINT NOT NULL,
    user_internal_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_comments_posts
        FOREIGN KEY (post_internal_id)
        REFERENCES posts(internal_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_comments_users
        FOREIGN KEY (user_internal_id)
        REFERENCES users(internal_id)
        ON DELETE CASCADE
)