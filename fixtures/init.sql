CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    bio VARCHAR,
    avatar_url VARCHAR UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX users_username_idx ON users (username);

CREATE TABLE polls (
    id INT PRIMARY KEY,
    creator_id INT REFERENCES users(id) ON DELETE SET NULL,
    topic VARCHAR NOT NULL,
    is_anonymous BOOLEAN NOT NULL,
    multiple_choice BOOLEAN NOT NULL,
    revote_ability BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE poll_options (
    id INT PRIMARY KEY,
    poll_id INT NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    index INT NOT NULL,
    option VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX poll_options_poll_id_idx ON poll_options (poll_id);

CREATE TABLE votes (
    id UUID PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    option_id INT NOT NULL REFERENCES poll_options(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX votes_user_id_option_id_idx ON votes (user_id, option_id);
