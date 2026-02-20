-- name: CreateRefreshToken :one
INSERT INTO refresh_token (
  token, user_id, expires_at
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: RevokeToken :exec
UPDATE refresh_token
SET expires_at = NOW(), updated_at = NOW()
WHERE token = $1;
