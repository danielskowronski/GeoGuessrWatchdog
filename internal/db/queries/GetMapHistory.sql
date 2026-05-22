-- name: GetMapHistory :many
SELECT
  mh.id,
  mh.bounds_min_lat,
  mh.bounds_min_lon,
  mh.bounds_max_lat,
  mh.bounds_max_lon,
  mh.max_err_distance,
  mh.api_updated,
  mh.location_count,
  mh.timestamp
FROM map_history as mh 
JOIN map_info as mi ON mh.map_id=mi.id
WHERE mi.api_id = $1
ORDER BY mh.timestamp ASC;