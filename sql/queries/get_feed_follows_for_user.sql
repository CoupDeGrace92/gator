-- name: GetFeedFollowsForUser :many
SELECT
    f.name AS feed_name,
    u.name AS user_name
FROM feed_follows
INNER JOIN users AS u
ON feed_follows.user_id = u.id
INNER JOIN feeds AS f
ON feed_follows.feed_id = f.id
WHERE u.name = $1;