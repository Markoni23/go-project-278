-- name: GetLinks :many
SELECT * FROM links
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetLinksCount :one
SELECT COUNT(1) FROM links;

-- name: GetLinkByShortName :one
SELECT * FROM links
WHERE short_name = $1
LIMIT 1;

-- name: GetLink :one
SELECT * FROM links
WHERE id = $1 LIMIT 1;

-- name: CreateLink :one
INSERT INTO links (
    original_url, short_name
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateLink :one
UPDATE links
    SET original_url = $2,
        short_name = $3
WHERE id = $1
RETURNING *;

-- name: DeleteLink :exec
DELETE FROM links
WHERE id = $1;