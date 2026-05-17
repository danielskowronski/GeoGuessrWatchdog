-- name: UpsertDivisionMapLastNotification :one
INSERT INTO division_map_last_notifications (
  division_name,
  game_mode,
  map_id,
  map_bounds_min_lat,
  map_bounds_min_lon,
  map_bounds_max_lat,
  map_bounds_max_lon,
  map_max_err_distance,
  map_api_updated,
  location_count
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (division_name, game_mode) DO UPDATE SET
  map_id = EXCLUDED.map_id,
  map_bounds_min_lat = EXCLUDED.map_bounds_min_lat,
  map_bounds_min_lon = EXCLUDED.map_bounds_min_lon,
  map_bounds_max_lat = EXCLUDED.map_bounds_max_lat,
  map_bounds_max_lon = EXCLUDED.map_bounds_max_lon,
  map_max_err_distance = EXCLUDED.map_max_err_distance,
  map_api_updated = EXCLUDED.map_api_updated,
  location_count = EXCLUDED.location_count,
  timestamp = now()
RETURNING id, timestamp;