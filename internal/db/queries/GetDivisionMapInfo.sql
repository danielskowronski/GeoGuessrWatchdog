-- name: GetDivisionMapInfo :one
SELECT
  division_name,
  game_mode,
  map_id,
  map_name,
  map_bounds_min_lat,
  map_bounds_min_lon,
  map_bounds_max_lat,
  map_bounds_max_lon,
  map_max_err_distance,
  map_api_updated,
  map_location_count
FROM division_map_info
WHERE division_name = $1
  AND game_mode = $2
LIMIT 1;