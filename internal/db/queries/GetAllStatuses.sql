-- name: GetAllStatuses :many
SELECT fetch_type,
  last_success
FROM fetch_status
WHERE last_success > now() - make_interval(days => @days);