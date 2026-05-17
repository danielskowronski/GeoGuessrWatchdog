-- name: GetDivisionMapLastNotification :one
SELECT
  id,
  map_id,
  map_bounds_min_lat,
  map_bounds_min_lon,
  map_bounds_max_lat,
  map_bounds_max_lon ,
  map_max_err_distance,
  map_api_updated,
  location_count,
  timestamp
FROM division_map_last_notifications
WHERE division_name = $1
  AND game_mode = $2
LIMIT 1;