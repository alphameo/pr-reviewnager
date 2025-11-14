CREATE TABLE user (
	id CHAR(36) PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	active BOOLEAN NOT NULL,
	UNIQUE(name)
);

CREATE TABLE team (
	id CHAR(36) PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	UNIQUE(name)
);

CREATE TABLE team_user (
	team_id CHAR(36) NOT NULL,
	user_id CHAR(36) NOT NULL,
	PRIMARY KEY (team_id, user_id),
	FOREIGN KEY (team_id) REFERENCES team(id)
		ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY (user_id) REFERENCES user(id)
		ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE pull_request (
	id CHAR(36) PRIMARY KEY,
	title VARCHAR(50) NOT NULL,
	author_id CHAR(36) NOT NULL,
	status ENUM('open', 'merged') NOT NULL,
	merged_at DATETIME,
	FOREIGN KEY (author_id) REFERENCES user(id)
		ON DELETE CASCADE ON UPDATE CASCADE,
	CHECK (
        (status = 'merged' AND merged_at IS NOT NULL) OR
        (status = 'open' AND merged_at IS NULL)
    )
);

CREATE TABLE pull_request_reviewer (
	pull_request_id char(36) NOT NULL,
	reviewer_id CHAR(36) NOT NULL,
	PRIMARY KEY (pull_request_id, reviewer_id),
	FOREIGN KEY (reviewer_id) REFERENCES user(id)
		ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY (pull_request_id) REFERENCES pull_request(id)
		ON DELETE CASCADE ON UPDATE CASCADE
);
