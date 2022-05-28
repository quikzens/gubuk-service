-- name: CreateTransaction :one
INSERT INTO transactions (
  id,
  tenant_id,
  owner_id,
  house_id,
  payment_status,
  payment_proof,
  total_payment,
  check_in,
  check_out,
  time_rent
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateTransactionStatusById :exec
UPDATE transactions 
SET 
  payment_status = $2,
  updated_at = $3
WHERE id = $1;

-- name: UpdateTransactionPaymentProofById :exec
UPDATE transactions 
SET 
  payment_status = $2,
  payment_proof = $3,
  updated_at = $4
WHERE id = $1;

-- name: DeleteTransaction :exec
DELETE FROM transactions 
WHERE id = $1;