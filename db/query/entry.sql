-- name: CreateEntry :one
INSERT INTO
  entries (account_id, amount)
VALUES
  (sqlc.arg('account_id'), sqlc.arg('amount')) RETURNING *;

-- name: GetEntry :one
SELECT
  *
FROM
  entries
WHERE
  id = sqlc.arg('id')
LIMIT
  1;

-- name: ListEntries :many
SELECT
  *
FROM
  entries
WHERE
  account_id = sqlc.arg('account_id')
ORDER BY
  id
LIMIT
  sqlc.arg('limit') OFFSET sqlc.arg('offset');