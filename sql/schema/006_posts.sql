-- +goose Up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url VARCHAR(150) NOT NULL, -- is this enough?
    description VARCHAR(250),
    published_at TIMESTAMP,
    feed_id UUID NOT NULL,
    UNIQUE(url)
);

-- +goose Down
DROP TABLE posts;