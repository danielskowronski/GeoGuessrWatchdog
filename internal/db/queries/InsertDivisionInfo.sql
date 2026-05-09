-- name: InsertDivisionInfo :one
INSERT INTO division_info (
  division_name,
  game_mode,
  map_id,
  map_name
) VALUES (
  $1, $2, $3, $4
)
RETURNING id;