CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Migration Up: create reactions table
DROP TABLE IF EXISTS reactions CASCADE;

CREATE TABLE IF NOT EXISTS reactions (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    user_internal_id BIGINT NOT NULL,
    post_internal_id BIGINT NOT NULL,
    reaction_type_internal_id BIGINT NOT NULL,
    emoji VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_reactions_user
        FOREIGN KEY (user_internal_id)
        REFERENCES users (internal_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_reactions_post
        FOREIGN KEY (post_internal_id)
        REFERENCES posts (internal_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_reactions_type
        FOREIGN KEY (reaction_type_internal_id)
        REFERENCES reaction_types (internal_id)
        ON DELETE CASCADE
);