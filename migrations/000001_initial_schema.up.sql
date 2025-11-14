CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    active BOOLEAN NOT NULL,
    UNIQUE(name)
);

CREATE TABLE team (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    UNIQUE(name)
);

CREATE TABLE team_user (
    team_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (team_id, user_id),
    FOREIGN KEY (team_id) REFERENCES team(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES "user"(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TYPE pull_request_status AS ENUM ('open', 'merged');

CREATE TABLE pull_request (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    author_id UUID NOT NULL,
    status pull_request_status NOT NULL,
    merged_at TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES "user"(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CHECK (
        (status = 'merged' AND merged_at IS NOT NULL) OR
        (status = 'open' AND merged_at IS NULL)
    )
);

CREATE TABLE pull_request_reviewer (
    pull_request_id UUID NOT NULL,
    reviewer_id UUID NOT NULL,
    PRIMARY KEY (pull_request_id, reviewer_id),
    FOREIGN KEY (reviewer_id) REFERENCES "user"(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (pull_request_id) REFERENCES pull_request(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);
