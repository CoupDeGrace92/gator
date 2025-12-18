-- name: GetFeeds :many
SELECT 
    f.name AS feed_name, 
    f.url AS feed_url, 
    u.name AS user
FROM feeds AS f
INNER JOIN users AS u
ON f.user_id = u.id;