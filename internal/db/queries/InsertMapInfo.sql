-- name: InsertMapInfo :one
INSERT INTO map_info (
  api_id,
  name,
  description,
  bounds_min_lat,
  bounds_min_lon,
  bounds_max_lat,
  bounds_max_lon,
  max_err_distance,
  api_updated,
  location_count
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id;