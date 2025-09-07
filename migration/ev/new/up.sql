-- Tabel untuk driver dengan struktur yang berbeda
CREATE TABLE drivers (
    driver_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(20) UNIQUE,
    date_of_birth DATE,
    residential_address JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabel untuk kendaraan dengan hubungan yang berbeda
CREATE TABLE vehicles (
    vehicle_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_plate VARCHAR(15) UNIQUE NOT NULL,
    vehicle_type VARCHAR(50),
    brand VARCHAR(50),
    model VARCHAR(50),
    manufacture_year INTEGER,
    max_capacity_kg DECIMAL(10,2),
    insurance_details JSONB,
    last_maintenance_date DATE
);

-- Tabel penghubung antara driver dan kendaraan (many-to-many)
CREATE TABLE driver_vehicle_assignments (
    assignment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id UUID REFERENCES drivers(driver_id),
    vehicle_id UUID REFERENCES vehicles(vehicle_id),
    assignment_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabel untuk shipments dengan struktur berbeda
CREATE TABLE shipments (
    shipment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assigned_driver_id UUID REFERENCES drivers(driver_id),
    assigned_vehicle_id UUID REFERENCES vehicles(vehicle_id),
    planned_delivery_date DATE,
    actual_delivery_date TIMESTAMP WITH TIME ZONE,
    pickup_address JSONB NOT NULL,
    delivery_address JSONB NOT NULL,
    current_status VARCHAR(20) DEFAULT 'CREATED',
    priority_level INTEGER DEFAULT 1,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabel untuk melacak history status shipment
CREATE TABLE shipment_status_history (
    status_history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID REFERENCES shipments(shipment_id),
    status VARCHAR(20) NOT NULL,
    status_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    location_coordinates VARCHAR,
    notes TEXT,
    updated_by UUID REFERENCES drivers(driver_id)
);

-- Tabel untuk package details (fitur baru)
CREATE TABLE shipment_packages (
    package_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID REFERENCES shipments(shipment_id),
    package_type VARCHAR(50),
    weight_kg DECIMAL(10,2),
    dimensions_cm JSONB, -- {length, width, height}
    description TEXT,
    is_fragile BOOLEAN DEFAULT FALSE,
    requires_signature BOOLEAN DEFAULT FALSE
);