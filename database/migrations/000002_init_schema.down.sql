DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS booking_status_histories;
DROP TABLE IF EXISTS booking_items;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS vehicle_images;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS vehicle_categories;
DROP TABLE IF EXISTS customer_documents;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;

-- DROP ENUMS
DROP TYPE IF EXISTS booking_source_enum;
DROP TYPE IF EXISTS customer_document_type_enum;
DROP TYPE IF EXISTS fuel_type_enum;
DROP TYPE IF EXISTS transmission_type_enum;
DROP TYPE IF EXISTS vehicle_status_enum;
DROP TYPE IF EXISTS booking_status_enum;
DROP TYPE IF EXISTS verification_status_enum;