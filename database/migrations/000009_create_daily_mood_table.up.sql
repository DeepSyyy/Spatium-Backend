CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS daily_moods (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    user_internal_id BIGINT NOT NULL,
    mood_tag_internal_id BIGINT NOT NULL,
    note TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_daily_moods_user FOREIGN KEY (user_internal_id)
        REFERENCES users (internal_id) ON DELETE CASCADE,

    CONSTRAINT fk_daily_moods_mood FOREIGN KEY (mood_tag_internal_id)
        REFERENCES mood_tags (internal_id) ON DELETE CASCADE,

    CONSTRAINT unique_user_date UNIQUE (user_internal_id, date)
);
