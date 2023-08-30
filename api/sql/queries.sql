-- name: CreateAccessInfo :exec
INSERT INTO store_log (product_id, "status") VALUES ($1, $2);

-- name: CheckStockExist :one
SELECT EXISTS(SELECT 1 FROM product_meta WHERE id = ($1));

-- name: GetStockMeta :one
SELECT id, "name", symbol, "description", "type", exchange,  "location"  FROM product_meta WHERE id = ($1);

-- name: GetAllStockMetaList :many
SELECT id, "name", symbol, "description", "type", exchange,  "location"  FROM product_meta;

-- name: GetStockMetaWithPlatform :one
SELECT product_meta.id, "name", symbol, "description", "type", exchange,  "location" , platform_name, identifier 
FROM product_meta 
JOIN product_platform 
ON product_meta.id = product_platform.product_id 
WHERE product_meta.id = ($1);

-- name: GetStockIdBySymbol :one
SELECT id FROM product_meta WHERE symbol = ($1);