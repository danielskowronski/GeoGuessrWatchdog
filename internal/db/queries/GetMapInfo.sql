-- name: GetMapInfo :one
SELECT
  id,
  api_id,
  name,
  description,
  bounds_min_lat,
  bounds_min_lon,
  bounds_max_lat,
  bounds_max_lon,
  max_err_distance,
  api_updated,
  location_count,
  last_changed
FROM map_info
WHERE api_id = $1
LIMIT 1;