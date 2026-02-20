-- name: CreateUser :one
INSERT INTO users (
  email, hashed_password
) VALUES ( $1, $2 )
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * 
FROM users
WHERE email = $1;

-- name: GetUserByRefreshToken :one
SELECT users.*, 
refresh_token.token as refresh_token, 
refresh_token.expires_at as expires_at
FROM users
INNER JOIN refresh_token
ON users.id = refresh_token.user_id
WHERE refresh_token.token = $1 
AND refresh_token.expires_at > NOW();
