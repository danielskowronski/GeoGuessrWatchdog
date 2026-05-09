-- name: InsertDivisionHistory :one
INSERT INTO division_history (
  division_id,
  map_id,
  map_name
) VALUES (
  $1, $2, $3
)
RETURNING id;