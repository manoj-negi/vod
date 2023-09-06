-- name: CreateVideo :one
INSERT INTO videos (
  name,
  uploaded_by
) VALUES (
  $1, $2
) RETURNING *;

-- name: ListVideos :many
SELECT $1 || name AS full_url, uploaded_by, is_active
FROM videos
ORDER BY name;