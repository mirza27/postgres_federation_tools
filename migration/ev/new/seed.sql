-- Drivers
INSERT INTO drivers (first_name, last_name, phone_number, date_of_birth, residential_address, is_active) VALUES
('James', 'Wilson', '+441234567890', '1985-03-15', '{"street": "12 Oxford Street", "city": "London", "postcode": "W1D 1BS"}', true),
('Sarah', 'Thompson', '+441234567891', '1990-08-22', '{"street": "45 Regent Street", "city": "London", "postcode": "W1B 2HE"}', true),
('Michael', 'Chen', '+441234567892', '1988-11-05', '{"street": "7 Covent Garden", "city": "London", "postcode": "WC2E 8RF"}', false);

-- Vehicles
INSERT INTO vehicles (license_plate, vehicle_type, brand, model, manufacture_year, max_capacity_kg, insurance_details) VALUES
('BD51 SMR', 'Van', 'Ford', 'Transit', 2020, 1500.00, '{"provider": "AXA", "policy_number": "AX123456"}'),
('LK06 EKY', 'Truck', 'Mercedes', 'Sprinter', 2019, 3500.00, '{"provider": "Allianz", "policy_number": "AL789012"}'),
('GX21 HSD', 'Refrigerated Van', 'Volkswagen', 'Crafter', 2021, 2000.00, '{"provider": "Aviva", "policy_number": "AV456789"}');

-- Assignments (Many-to-Many)
INSERT INTO driver_vehicle_assignments (driver_id, vehicle_id, assignment_date) VALUES
((SELECT driver_id FROM drivers WHERE phone_number = '+441234567890'), (SELECT vehicle_id FROM vehicles WHERE license_plate = 'BD51 SMR'), '2023-10-01'),
((SELECT driver_id FROM drivers WHERE phone_number = '+441234567891'), (SELECT vehicle_id FROM vehicles WHERE license_plate = 'LK06 EKY'), '2023-10-02'),
((SELECT driver_id FROM drivers WHERE phone_number = '+441234567892'), (SELECT vehicle_id FROM vehicles WHERE license_plate = 'GX21 HSD'), '2023-10-03');

-- Shipments
INSERT INTO shipments (assigned_driver_id, assigned_vehicle_id, planned_delivery_date, pickup_address, delivery_address, current_status) VALUES
((SELECT driver_id FROM drivers WHERE phone_number = '+441234567890'), (SELECT vehicle_id FROM vehicles WHERE license_plate = 'BD51 SMR'), '2023-10-05', 
 '{"street": "32 Baker Street", "city": "London", "postcode": "NW1 6XE"}', 
 '{"street": "15 Victoria Road", "city": "Manchester", "postcode": "M4 4AE"}', 'IN_TRANSIT'),
((SELECT driver_id FROM drivers WHERE phone_number = '+441234567891'), (SELECT vehicle_id FROM vehicles WHERE license_plate = 'LK06 EKY'), '2023-10-06',
 '{"street": "67 Kingsway", "city": "London", "postcode": "WC2B 6TD"}', 
 '{"street": "22 Railway Street", "city": "Liverpool", "postcode": "L1 0BS"}', 'DELIVERED');

-- Shipment Packages
INSERT INTO shipment_packages (shipment_id, package_type, weight_kg, dimensions_cm, is_fragile) VALUES
((SELECT shipment_id FROM shipments LIMIT 1 OFFSET 0), 'Electronics', 5.50, '{"length": 40, "width": 30, "height": 20}', true),
((SELECT shipment_id FROM shipments LIMIT 1 OFFSET 1), 'Furniture', 120.75, '{"length": 150, "width": 80, "height": 50}', false);

-- Status History
INSERT INTO shipment_status_history (shipment_id, status, location_coordinates, notes) VALUES
((SELECT shipment_id FROM shipments LIMIT 1 OFFSET 0), 'PICKED_UP', '51.5072,-0.1276', 'Picked up from warehouse'),
((SELECT shipment_id FROM shipments LIMIT 1 OFFSET 0), 'IN_TRANSIT', '53.4808,-2.2426', 'On the way to destination');