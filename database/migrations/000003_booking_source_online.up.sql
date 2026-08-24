ALTER TYPE booking_source_enum ADD VALUE 'online';
ALTER TABLE bookings ALTER COLUMN booking_source SET DEFAULT 'online';
