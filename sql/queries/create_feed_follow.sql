-- name: CreateFeedFollows :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES(
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT 
    i.*,
    feeds.name AS feed_name,
    users.name AS user
FROM inserted_feed_follow AS i
INNER JOIN feeds
ON i.feed_id = feeds.id
INNER JOIN users
ON i.user_id = users.id;