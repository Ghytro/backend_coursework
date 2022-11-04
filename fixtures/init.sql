CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    first_name VARCHAR,
    last_name VARCHAR,
    password VARCHAR NOT NULL,
    bio VARCHAR,
    avatar_url VARCHAR UNIQUE,
    country VARCHAR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX users_username_idx ON users (username);

INSERT INTO users
    (username, first_name, last_name, password, bio, country)
VALUES
    ('ghytro', 'Michael', 'Korobkov', crypt('123123', gen_salt('bf')), 'some bio that i want to tell you', 'RU');

CREATE TABLE polls (
    id SERIAL PRIMARY KEY,
    creator_id INT REFERENCES users(id) ON DELETE SET NULL,
    topic VARCHAR NOT NULL,
    is_anonymous BOOLEAN NOT NULL DEFAULT FALSE,
    multiple_choice BOOLEAN NOT NULL DEFAULT FALSE,
    revote_ability BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE poll_options (
    id SERIAL PRIMARY KEY,
    poll_id INT NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    index INT NOT NULL,
    option VARCHAR NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX poll_options_poll_id_idx ON poll_options (poll_id);

CREATE TABLE votes (
    id UUID PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    poll_id INT NOT NULL REFERENCES polls(id) ON DELETE CASCADE, -- для быстроты получения
    option_id INT NOT NULL REFERENCES poll_options(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX votes_user_id_option_id_idx ON votes (user_id, option_id);
