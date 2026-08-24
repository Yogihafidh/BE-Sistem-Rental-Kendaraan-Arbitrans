-- name: CheckCustomer :one
SELECT id, verification_status
FROM customers
WHERE deleted_at IS NULL
  AND (nik = $1 OR phone = $2)
LIMIT 1;

-- name: ListCustomers :many
SELECT * FROM customers
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateCustomer :one
INSERT INTO customers (full_name, email, phone, nik, address, date_of_birth)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateCustomer :one
UPDATE customers
SET full_name = $2,
    email = $3,
    phone = $4,
    nik = $5,
    address = $6,
    date_of_birth = $7,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateCustomerVerificationStatus :one
UPDATE customers
SET verification_status = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCustomer :exec
UPDATE customers
SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateCustomerDocument :one
INSERT INTO customer_documents (customer_id, document_type, file_url)
VALUES ($1, $2, $3)
ON CONFLICT (customer_id, document_type)
DO UPDATE SET file_url = EXCLUDED.file_url, uploaded_at = now(), verified_by = NULL, verified_at = NULL
RETURNING *;

-- name: ListCustomerDocuments :many
SELECT * FROM customer_documents
WHERE customer_id = $1
ORDER BY uploaded_at DESC;

-- name: DeleteCustomerDocument :exec
DELETE FROM customer_documents
WHERE id = $1;

-- name: VerifyCustomerDocument :one
UPDATE customer_documents
SET verified_by = $2, verified_at = now()
WHERE id = $1
RETURNING *;
