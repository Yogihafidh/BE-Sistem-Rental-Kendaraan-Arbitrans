-- ENUM
CREATE TYPE verification_status_enum AS ENUM ('pending', 'verified', 'rejected');
CREATE TYPE booking_status_enum AS ENUM ('pending', 'confirmed', 'ongoing', 'cancelled', 'completed');
CREATE TYPE vehicle_status_enum AS ENUM ('available', 'booked', 'maintenance', 'inactive');
CREATE TYPE transmission_type_enum AS ENUM ('manual', 'automatic');

CREATE TYPE fuel_type_enum AS ENUM ('gasoline', 'diesel', 'electric', 'hybrid');
CREATE TYPE customer_document_type_enum AS ENUM ('ktp-penyewa', 'ktp-penjamin', 'sim', 'passport', 'tiket', 'npwp');

CREATE TYPE booking_source_enum AS ENUM ('website', 'mobile_app', 'call_center', 'walk_in');

-- SCHEMA TABLES
CREATE TABLE roles(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE CHECK (char_length(trim(name)) > 0),
    description TEXT
);

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key to roles (many-to-one)
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,

    full_name VARCHAR(255) NOT NULL CHECK (char_length(trim(full_name)) >= 3),
    email VARCHAR(255) NOT NULL UNIQUE CHECK (position('@' in email) > 1),
    phone VARCHAR(20) CHECK (phone IS NULL OR char_length(phone) BETWEEN 10 AND 20),
    password_hash TEXT NOT NULL CHECK (char_length(password_hash) > 20),

    is_active BOOLEAN NOT NULL DEFAULT true,

    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL CHECK (char_length(trim(full_name)) >= 3),
    email VARCHAR(255) UNIQUE CHECK (position('@' in email) > 1),
    phone VARCHAR(20) NOT NULL CHECK (phone ~ '^[0-9+]+$'),
    nik VARCHAR(50) UNIQUE NOT NULL CHECK (nik ~ '^[0-9]{16}$'),
    address TEXT,
    date_of_birth DATE,
    verification_status verification_status_enum NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE customer_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Foreign key to customers (one-to-many)
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    document_type customer_document_type_enum NOT NULL,
    file_url TEXT NOT NULL CHECK (char_length(trim(file_url)) > 0),
    -- Foreign key to users (many-to-one) for tracking who verified the document
    verified_by UUID REFERENCES users(id) ON DELETE SET NULL,
    verified_at TIMESTAMPTZ,
    uploaded_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (customer_id, document_type)
);

CREATE TABLE vehicle_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE CHECK (char_length(trim(name)) > 0),
    description TEXT
);

CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Foreign key to vehicle_categories (many-to-one)
    category_id UUID NOT NULL REFERENCES vehicle_categories(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year SMALLINT NOT NULL CHECK ( year >= 1900 AND year <= EXTRACT(YEAR FROM now()) + 1),
    transmission transmission_type_enum NOT NULL,
    fuel_type fuel_type_enum NOT NULL,
    seat_capacity INT CHECK (seat_capacity > 0),
    luggage_capacity INT CHECK (luggage_capacity >= 0),
    plate_number VARCHAR(50) NOT NULL UNIQUE CHECK (char_length(trim(plate_number)) >= 3),
    daily_price NUMERIC(12,2) NOT NULL CHECK (daily_price >= 0),
    status vehicle_status_enum NOT NULL DEFAULT 'available',
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE vehicle_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Foreign key to vehicles (many-to-one)
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL CHECK (char_length(trim(image_url)) > 0),
    is_thumbnail BOOLEAN DEFAULT false,
    uploaded_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key to customers (many-to-one)
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,

    -- Foreign key to users (many-to-one) for tracking who created the booking
    created_by_user UUID REFERENCES users(id) ON DELETE SET NULL,

    booking_code VARCHAR(100) NOT NULL UNIQUE CHECK (char_length(trim(booking_code)) > 0),
    booking_source booking_source_enum NOT NULL DEFAULT 'website',
    pickup_location TEXT NOT NULL,
    return_location TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    pickup_time TIME,
    return_time TIME,

    booking_status booking_status_enum NOT NULL DEFAULT 'pending',

    subtotal NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
    total_amount NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CHECK (end_date >= start_date)
);

CREATE TABLE booking_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key to bookings (many-to-one)
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    -- Foreign key to vehicles (many-to-one)
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE RESTRICT,

    price_per_day NUMERIC(12,2) NOT NULL CHECK (price_per_day >= 0),
    total_days INT NOT NULL CHECK (total_days > 0),
    subtotal NUMERIC(12,2) NOT NULL CHECK (subtotal >= 0),

    UNIQUE (booking_id, vehicle_id)
);

CREATE TABLE booking_status_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key to bookings (many-to-one)
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,

    -- Foreign key to users (many-to-one) for tracking who changed the status
    changed_by UUID REFERENCES users(id) ON DELETE SET NULL,

    status booking_status_enum NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key to users (many-to-one) for tracking who performed the action
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMPTZ DEFAULT now()
);

