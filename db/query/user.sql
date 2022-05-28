-- name: CreateUser :one
INSERT INTO users (
  id,
  fullname,
  username,
  email,
  role,
  gender,
  phone_number,
  password,
  address,
  avatar 
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id;

-- name: UpdateUserById :exec
UPDATE users
SET
  fullname = $2,
  email = $3,
  gender = $4,
  phone_number = $5,
  address = $6,
  updated_at = $7
WHERE id = $1;

-- name: UpdateUserAvatarById :exec
UPDATE users 
SET 
  avatar = $2,
  updated_at = $3
WHERE id = $1;

-- name: UpdateUserPasswordById :exec
UPDATE users 
SET 
  password = $2,
  updated_at = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT 
  id, 
  fullname, 
  username,
  email,
  role,
  gender,
  phone_number,
  password,
  address, 
  avatar,
  created_at, 
  updated_at
FROM users
WHERE users.username = $1 LIMIT 1;

-- name: GetUserById :one
SELECT 
  id, 
  fullname, 
  username,
  email,
  role,
  gender,
  phone_number,
  address, 
  avatar,
  created_at, 
  updated_at
FROM users
WHERE users.id = $1 LIMIT 1;

-- name: GetUserAvatarById :one
SELECT avatar FROM users WHERE users.id = $1 LIMIT 1;