-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
  sqlc.arg('owner'), sqlc.arg('balance'), sqlc.arg('currency')
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = sqlc.arg('id')
LIMIT sqlc.arg('limit');

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: UpdateAccount :one
UPDATE accounts
SET balance = sqlc.arg('balance')
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = sqlc.arg('id');