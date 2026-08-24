-- name: GetDashboardSummary :one
SELECT
    (SELECT COUNT(*) FROM vehicles WHERE deleted_at IS NULL) AS total_vehicles,
    (SELECT COUNT(*) FROM vehicles WHERE deleted_at IS NULL AND status = 'available') AS available_vehicles,
    (SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL) AS total_customers,
    (SELECT COUNT(*) FROM bookings) AS total_bookings,
    (SELECT COUNT(*) FROM bookings WHERE booking_status = 'pending') AS pending_bookings,
    (SELECT COUNT(*) FROM bookings WHERE booking_status = 'confirmed') AS confirmed_bookings,
    (SELECT COUNT(*) FROM bookings WHERE booking_status = 'completed') AS completed_bookings,
    COALESCE((SELECT SUM(total_amount) FROM bookings WHERE booking_status = 'completed'), 0)::numeric AS total_revenue;
