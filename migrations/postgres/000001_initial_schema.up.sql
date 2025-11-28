-- +migrate Up

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS team (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS team_user (
    team_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (team_id, user_id),
    FOREIGN KEY (team_id) REFERENCES team (id)
    ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
    ON DELETE CASCADE ON UPDATE CASCADE
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pull_request_status') THEN
        CREATE TYPE pull_request_status AS ENUM ('open', 'merged');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS pull_request (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR NOT NULL,
    author_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    status PULL_REQUEST_STATUS NOT NULL,
    merged_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (author_id) REFERENCES "user" (id)
    ON DELETE CASCADE ON UPDATE CASCADE,
    CHECK (
        (status = 'merged' AND merged_at IS NOT NULL)
        OR (status = 'open' AND merged_at IS NULL)
    )
);

CREATE TABLE IF NOT EXISTS pull_request_reviewer (
    pull_request_id UUID NOT NULL,
    reviewer_id UUID NOT NULL,
    PRIMARY KEY (pull_request_id, reviewer_id),
    FOREIGN KEY (reviewer_id) REFERENCES "user" (id)
    ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (pull_request_id) REFERENCES pull_request (id)
    ON DELETE CASCADE ON UPDATE CASCADE
);
