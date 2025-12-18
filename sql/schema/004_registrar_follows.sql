-- +goose Up
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
SELECT
    gen_random_uuid(),
    NOW(),
    NOW(),
    f.user_id,
    f.id
FROM feeds as f
ON CONFLICT (user_id, feed_id) DO NOTHING;


-- +goose Down
DELETE FROM feed_follows
WHERE (user_id, feed_id) IN (
    SELECT
        f.user_id,
        f.id
    FROM feeds AS f
);
