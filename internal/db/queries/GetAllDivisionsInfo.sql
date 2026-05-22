-- name: GetAllDivisionsInfo :many
SELECT
    id,
    division_name,
    game_mode,
    map_id,
    map_name,
    last_changed
FROM division_info
ORDER BY id ASC;