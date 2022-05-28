-- name: CreateHouse :one
INSERT INTO homes (
  id,
  owner_id,
  title,
  featured_image,
  bedrooms,
  bathrooms,
  type_rent,
  price,
  province_id,
  city_id,
  description,
  amenities,
  area
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: UpdateHouse :one
UPDATE homes
SET
  title = $2,
  featured_image = $3,
  bedrooms = $4,
  bathrooms = $5,
  type_rent = $6,
  price = $7,
  province_id = $8,
  city_id = $9,
  description = $10,
  amenities = $11,
  area = $12
WHERE id = $1
RETURNING *;

-- name: DeleteHouse :exec
DELETE FROM homes 
WHERE id = $1;

-- name: GetHouseById :one
SELECT 
  homes.id,
  homes.title,
  homes.featured_image,
  homes.bedrooms,
  homes.bathrooms,
  homes.type_rent,
  homes.price,
  homes.province_id,
  homes.city_id,
  homes.description,
  homes.amenities,
  homes.area,
  homes.created_at,
  homes.updated_at,
  owner.id AS owner_id,
  owner.fullname AS owner_fullname,
  owner.username AS owner_username,
  owner.email AS owner_email,
  owner.role AS owner_role, 
  owner.gender AS owner_gender, 
  owner.phone_number AS owner_phone_number, 
  owner.address AS owner_address
FROM homes
JOIN users AS owner
ON owner.id = homes.owner_id
WHERE homes.id = $1 LIMIT 1;

-- name: ListHouse :many
SELECT 
  id,
  title,
  featured_image,
  bedrooms,
  bathrooms,
  type_rent,
  price,
  province_id,
  city_id,
  description,
  amenities,
  area,
  created_at,
  updated_at
FROM homes 
ORDER BY created_at DESC;

-- name: ListMyHouse :many
SELECT 
  id,
  title,
  featured_image,
  bedrooms,
  bathrooms,
  type_rent,
  price,
  province_id,
  city_id,
  description,
  amenities,
  area,
  created_at,
  updated_at
FROM homes 
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: CountHouse :one
SELECT COUNT(*) FROM homes;