-- +goose Up
ALTER TABLE feeds
ADD last_fetched_at TIMESTAMP;


-- +goose Down
ALTER TABLE
DROP COLUMN last_fetched_at;