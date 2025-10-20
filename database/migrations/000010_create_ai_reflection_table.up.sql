CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS ai_reflections (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    user_internal_id BIGINT NOT NULL,
    mood_tag_internal_id BIGINT NOT NULL,
    reflection TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_ai_reflections_user
        FOREIGN KEY (user_internal_id)
        REFERENCES users (internal_id)
        ON DELETE CASCADE
);
