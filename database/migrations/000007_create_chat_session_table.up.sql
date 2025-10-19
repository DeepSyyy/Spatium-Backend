CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS chat_sessions (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    user_internal_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    mood_tag VARCHAR(50) DEFAULT 'neutral',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_chat_sessions_user
        FOREIGN KEY (user_internal_id)
        REFERENCES users (internal_id)
        ON DELETE CASCADE
);