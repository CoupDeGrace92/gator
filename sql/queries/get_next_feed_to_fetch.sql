-- name: GetNextFeedToFetch :one
SELECT url FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;