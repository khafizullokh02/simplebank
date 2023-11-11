-- name: CreateTransfer :one
INSERT INTO
  transfers (from_account_id, to_account_id, amount)
VALUES
  (
    sqlc.arg('from_account_id'),
    sqlc.arg('to_account_id'),
    sqlc.arg('amount')
  ) RETURNING *;

-- name: GetTransfer :one
SELECT
  *
FROM
  transfers
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListTransfers :many
SELECT
  *
FROM
  transfers
WHERE
  from_account_id = sqlc.arg('from_account_id')
  OR to_account_id = sqlc.arg('to_account_id')
ORDER BY
  id
LIMIT
  sqlc.arg('limit') OFFSET sqlc.arg('offset');