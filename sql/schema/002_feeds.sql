-- +goose Up
CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    user_id uuid NOT NULL,
    FOREIGN KEY (user_id) references users(id)
    ON DELETE CASCADE);

-- +goose Down
DROP TABLE feeds;