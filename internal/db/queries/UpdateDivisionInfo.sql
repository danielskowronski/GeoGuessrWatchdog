-- name: UpdateDivisionInfo :one
UPDATE division_info
SET
  map_id = $2,
  map_name = $3
WHERE id = $1
RETURNING id;