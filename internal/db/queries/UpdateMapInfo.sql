-- name: UpdateMapInfo :one
UPDATE map_info
SET
  name = $2,
  description = $3,
  bounds_min_lat = $4,
  bounds_min_lon = $5,
  bounds_max_lat = $6,
  bounds_max_lon = $7,
  max_err_distance = $8,
  api_updated = $9,
  location_count = $10
WHERE api_id = $1
RETURNING id;