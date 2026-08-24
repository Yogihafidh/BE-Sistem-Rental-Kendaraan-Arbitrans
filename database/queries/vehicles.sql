-- name: ListVehicleCategories :many
SELECT * FROM vehicle_categories
ORDER BY name ASC;

-- name: CreateVehicleCategory :one
INSERT INTO vehicle_categories (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetVehicleCategoryByID :one
SELECT * FROM vehicle_categories
WHERE id = $1
LIMIT 1;

-- name: UpdateVehicleCategory :one
UPDATE vehicle_categories
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteVehicleCategory :exec
DELETE FROM vehicle_categories
WHERE id = $1;

-- name: ListVehicles :many
SELECT * FROM vehicles
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateVehicle :one
INSERT INTO vehicles (
    category_id, name, brand, model, year, transmission, fuel_type,
    seat_capacity, luggage_capacity, plate_number, daily_price, status, description
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetVehicleByID :one
SELECT * FROM vehicles
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateVehicle :one
UPDATE vehicles
SET category_id = $2,
    name = $3,
    brand = $4,
    model = $5,
    year = $6,
    transmission = $7,
    fuel_type = $8,
    seat_capacity = $9,
    luggage_capacity = $10,
    plate_number = $11,
    daily_price = $12,
    description = $13,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateVehicleStatus :one
UPDATE vehicles
SET status = $2, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteVehicle :exec
UPDATE vehicles
SET deleted_at = now(), updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAvailableVehicles :many
SELECT v.*
FROM vehicles v
WHERE v.deleted_at IS NULL
  AND v.status = 'available'
  AND NOT EXISTS (
    SELECT 1
    FROM booking_items bi
    JOIN bookings b ON b.id = bi.booking_id
    WHERE bi.vehicle_id = v.id
      AND b.booking_status IN ('pending', 'confirmed', 'ongoing')
      AND b.start_date <= $2
      AND b.end_date >= $1
  )
ORDER BY v.created_at DESC;

-- name: CountUnavailableVehicleBookings :one
SELECT COUNT(*)
FROM booking_items bi
JOIN bookings b ON b.id = bi.booking_id
WHERE bi.vehicle_id = $1
  AND b.booking_status IN ('pending', 'confirmed', 'ongoing')
  AND b.start_date <= $3
  AND b.end_date >= $2;

-- name: CreateVehicleImage :one
INSERT INTO vehicle_images (vehicle_id, image_url, is_thumbnail)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListVehicleImages :many
SELECT * FROM vehicle_images
WHERE vehicle_id = $1
ORDER BY uploaded_at DESC;

-- name: DeleteVehicleImage :exec
DELETE FROM vehicle_images
WHERE id = $1;

-- name: UnsetVehicleThumbnails :exec
UPDATE vehicle_images
SET is_thumbnail = false
WHERE vehicle_id = $1;

-- name: SetVehicleImageThumbnail :one
UPDATE vehicle_images
SET is_thumbnail = true
WHERE id = $1
RETURNING *;

-- name: GetVehicleImageByID :one
SELECT * FROM vehicle_images
WHERE id = $1
LIMIT 1;
