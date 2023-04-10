-- name: CreateAccessInfo :exec
INSERT INTO access_log (connected_at) VALUES (NOW());