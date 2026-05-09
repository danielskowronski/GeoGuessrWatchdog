-- name: GetDivisionInfo :one
SELECT
    id,
    division_name,
    game_mode,
    map_id,
    map_name,
    last_changed
FROM division_info
WHERE division_name = $1
  AND game_mode = $2
LIMIT 1;