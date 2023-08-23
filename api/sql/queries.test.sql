-- name: CountTestTableEntity :one
SELECT COUNT(*) FROM test_table;

-- name: InsertTestTableEntity :exec
INSERT INTO test_table (id) VALUES (DEFAULT);