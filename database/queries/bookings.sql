-- name: CreateBooking :one
INSERT INTO bookings (
    customer_id, created_by_user, booking_code, booking_source,
    pickup_location, return_location, start_date, end_date,
    pickup_time, return_time, booking_status, subtotal, total_amount, notes
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: ListBookings :many
SELECT * FROM bookings
ORDER BY created_at DESC;

-- name: GetBookingByID :one
SELECT * FROM bookings
WHERE id = $1
LIMIT 1;

-- name: GetBookingByCode :one
SELECT * FROM bookings
WHERE booking_code = $1
LIMIT 1;

-- name: UpdateBooking :one
UPDATE bookings
SET pickup_location = $2,
    return_location = $3,
    start_date = $4,
    end_date = $5,
    pickup_time = $6,
    return_time = $7,
    subtotal = $8,
    total_amount = $9,
    notes = $10,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateBookingStatus :one
UPDATE bookings
SET booking_status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBooking :exec
DELETE FROM bookings
WHERE id = $1;

-- name: CreateBookingItem :one
INSERT INTO booking_items (booking_id, vehicle_id, price_per_day, total_days, subtotal)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListBookingItems :many
SELECT * FROM booking_items
WHERE booking_id = $1
ORDER BY id ASC;

-- name: DeleteBookingItem :exec
DELETE FROM booking_items
WHERE id = $1;

-- name: CreateBookingStatusHistory :one
INSERT INTO booking_status_histories (booking_id, changed_by, status, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListBookingStatusHistories :many
SELECT * FROM booking_status_histories
WHERE booking_id = $1
ORDER BY created_at ASC;

-- name: CreateAuditLog :one
INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_value, new_value)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
