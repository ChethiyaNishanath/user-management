-- name: CreateUser :one
INSERT INTO USERS (USER_ID, FIRST_NAME, LAST_NAME, EMAIL, PHONE, AGE, STATUS)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: FindUserById :one
SELECT * FROM USERS WHERE USER_ID = $1 LIMIT 1;

-- name: ListAllUsersPaged :many
SELECT * FROM USERS LIMIT $1 OFFSET $2;

-- name: DeleteUserByID :exec
DELETE FROM USERS WHERE USER_ID = $1;

-- name: UpdateUser :one
UPDATE users
SET
    FIRST_NAME = COALESCE(sqlc.narg('first_name'), FIRST_NAME),
    LAST_NAME  = COALESCE(sqlc.narg('last_name'), LAST_NAME),
    EMAIL      = COALESCE(sqlc.narg('email'), EMAIL),
    PHONE      = COALESCE(sqlc.narg('phone'), PHONE),
    AGE        = COALESCE(sqlc.narg('age'), AGE),
    STATUS     = COALESCE(sqlc.narg('status'), STATUS)
WHERE user_id = sqlc.arg('user_id')
RETURNING *;