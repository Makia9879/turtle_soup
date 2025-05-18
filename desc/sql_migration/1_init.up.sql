-- 创建活动token表
CREATE TABLE IF NOT EXISTS T_activity_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建会话token表
CREATE TABLE IF NOT EXISTS T_session_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(64) NOT NULL UNIQUE,
    activity_token_id INTEGER REFERENCES activity_tokens(id),
    remaining_attempts INTEGER NOT NULL DEFAULT 3,
    remaining_answers INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建海龟汤故事表
CREATE TABLE IF NOT EXISTS T_turtle_soup_stories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    answer TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户会话关联表
CREATE TABLE IF NOT EXISTS T_user_sessions (
    id SERIAL PRIMARY KEY,
    session_token_id INTEGER REFERENCES session_tokens(id),
    story_id INTEGER REFERENCES turtle_soup_stories(id),
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);