-- name: CreateInstrument :one
INSERT INTO INSTRUMENTS (ID, SYMBOL, NAME, INSTRUMENT_TYPE, EXCHANGE, LAST_PRICE, CREATED_AT, UPDATED_AT)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: FindInstrumentById :one
SELECT * FROM INSTRUMENTS WHERE ID = $1 LIMIT 1;

-- name: ListAllInstrumentPaged :many
SELECT * FROM INSTRUMENTS LIMIT $1 OFFSET $2;

-- name: DeleteInstrumentById :exec
DELETE FROM INSTRUMENTS WHERE ID = $1;

-- name: UpdateInstrument :one
UPDATE INSTRUMENTS
SET
    SYMBOL = COALESCE(sqlc.narg('symbol'), SYMBOL),
    NAME  = COALESCE(sqlc.narg('name'), NAME),
    INSTRUMENT_TYPE      = COALESCE(sqlc.narg('instrument_type'), EXCHANGE),
    EXCHANGE      = COALESCE(sqlc.narg('exchange'), PHONE),
    LAST_PRICE        = COALESCE(sqlc.narg('last_price'), LAST_PRICE),
    CREATED_AT     = COALESCE(sqlc.narg('created_at'), CREATED_AT),
    UPDATED_AT     = COALESCE(sqlc.narg('updated_at'), UPDATED_AT)
WHERE ID = sqlc.arg('id')
RETURNING *;