-- name: GetPostsFromUser :many
SELECT
    f.name AS feed_name,
    p.title AS title,
    p.url AS url
FROM posts AS p
INNER JOIN feeds AS f
ON f.id = p.feed_id
INNER JOIN feed_follows AS fol
ON p.feed_id = fol.feed_id
WHERE fol.user_id = $1
ORDER BY published_at DESC
LIMIT $2;