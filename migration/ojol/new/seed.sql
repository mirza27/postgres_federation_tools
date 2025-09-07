-- Insert data untuk tabel users (10 data: 5 customer + 5 driver)
INSERT INTO users (user_id, email, username, password, user_type, status) VALUES
-- Customer accounts
('10000000-0000-0000-0000-000000000001', 'customer1@mail.com', 'cust1', 'pass1', 'customer', 'active'),
('10000000-0000-0000-0000-000000000002', 'customer2@mail.com', 'cust2', 'pass2', 'customer', 'active'),
('10000000-0000-0000-0000-000000000003', 'customer3@mail.com', 'cust3', 'pass3', 'customer', 'active'),
('10000000-0000-0000-0000-000000000004', 'customer4@mail.com', 'cust4', 'pass4', 'customer', 'inactive'),
('10000000-0000-0000-0000-000000000005', 'customer5@mail.com', 'cust5', 'pass5', 'customer', 'active'),
-- Driver accounts
('20000000-0000-0000-0000-000000000001', 'driver1@mail.com', 'driv1', 'pass6', 'driver', 'active'),
('20000000-0000-0000-0000-000000000002', 'driver2@mail.com', 'driv2', 'pass7', 'driver', 'active'),
('20000000-0000-0000-0000-000000000003', 'driver3@mail.com', 'driv3', 'pass8', 'driver', 'inactive'),
('20000000-0000-0000-0000-000000000004', 'driver4@mail.com', 'driv4', 'pass9', 'driver', 'active'),
('20000000-0000-0000-0000-000000000005', 'driver5@mail.com', 'driv5', 'pass10', 'driver', 'active');

-- Insert data untuk tabel customers
INSERT INTO customers (customer_id, user_id, name, phone, id_card, status, address) VALUES
('30000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'Budi Santoso', '+62811111111', '1234567890', 'active', 'Jl. Merdeka No. 1, Jakarta'),
('30000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', 'Sari Dewi', '+62812222222', '2345678901', 'active', 'Jl. Sudirman No. 12, Jakarta'),
('30000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000003', 'Ahmad Fauzi', '+62813333333', '3456789012', 'active', 'Jl. Thamrin No. 23, Jakarta'),
('30000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000004', 'Diana Putri', '+62814444444', '4567890123', 'inactive', 'Jl. Gatot Subroto No. 45, Jakarta'),
('30000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000005', 'Rizal Hermawan', '+62815555555', '5678901234', 'active', 'Jl. Asia Afrika No. 56, Bandung');

-- Insert data untuk tabel drivers
INSERT INTO drivers (driver_id, user_id, name, type, phone, id_card, status, vehicle_info) VALUES
('40000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'Joko Prasetyo', 'car', '+62816666666', '6789012345', 'active', '{"plate": "B 1234 ABC", "type": "SUV", "color": "Black"}'),
('40000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000002', 'Dewi Anggraini', 'motorcycle', '+62817777777', '7890123456', 'active', '{"plate": "B 5678 DEF", "type": "Sport", "color": "Red"}'),
('40000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000003', 'Rudi Hartono', 'car', '+62818888888', '8901234567', 'inactive', '{"plate": "B 9012 GHI", "type": "MPV", "color": "White"}'),
('40000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000004', 'Maya Sari', 'car', '+62819999999', '9012345678', 'active', '{"plate": "B 3456 JKL", "type": "Hatchback", "color": "Blue"}'),
('40000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000005', 'Fajar Siddiq', 'motorcycle', '+62810000000', '0123456789', 'active', '{"plate": "B 7890 MNO", "type": "Matic", "color": "Green"}');

-- Insert data untuk tabel orders
INSERT INTO orders (order_id, customer_id, order_type, status, amount, note, depart_location, purpose_location) VALUES
('50000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'instant', 'completed', 25000.500, 'Jangan terlambat', '{"address": "Jl. Merdeka No. 1, Jakarta", "lat": -6.200, "lon": 106.816}', '{"address": "Plaza Indonesia, Jakarta", "lat": -6.193, "lon": 106.823}'),
('50000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', 'scheduled', 'processing', 45000.000, 'Bawa helm extra', '{"address": "Jl. Sudirman No. 12, Jakarta", "lat": -6.208, "lon": 106.818}', '{"address": "GBK Arena, Jakarta", "lat": -6.219, "lon": 106.803}'),
('50000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', 'instant', 'cancelled', 30000.750, 'Penumpang 2 orang', '{"address": "Jl. Thamrin No. 23, Jakarta", "lat": -6.186, "lon": 106.823}', '{"address": "Monas, Jakarta", "lat": -6.175, "lon": 106.827}'),
('50000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000001', 'scheduled', 'completed', 60000.000, 'Tidak ada catatan', '{"address": "Plaza Indonesia, Jakarta", "lat": -6.193, "lon": 106.823}', '{"address": "Bundaran HI, Jakarta", "lat": -6.194, "lon": 106.823}'),
('50000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000005', 'instant', 'processing', 35000.000, 'Ambil di gerbang', '{"address": "Jl. Asia Afrika No. 56, Bandung", "lat": -6.921, "lon": 107.610}', '{"address": "Stasiun Bandung, Bandung", "lat": -6.917, "lon": 107.618}');

-- Insert data untuk tabel assignments
INSERT INTO assignments (assignment_id, order_id, driver_id, status) VALUES
('60000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000001', 'accepted'),
('60000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', '40000000-0000-0000-0000-000000000002', 'accepted'),
('60000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000003', '40000000-0000-0000-0000-000000000003', 'rejected'),
('60000000-0000-0000-0000-000000000004', '50000000-0000-0000-0000-000000000004', '40000000-0000-0000-0000-000000000004', 'accepted'),
('60000000-0000-0000-0000-000000000005', '50000000-0000-0000-0000-000000000005', '40000000-0000-0000-0000-000000000005', 'pending');

-- Insert data untuk tabel reviews
INSERT INTO reviews (review_id, order_id, customer_id, review, star) VALUES
('70000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'Pelayanan bagus', 5),
('70000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', 'Sopirnya ramah', 4),
('70000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000001', 'Kendaraan nyaman', 5),
('70000000-0000-0000-0000-000000000004', '50000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000005', 'Tepat waktu', 4),
('70000000-0000-0000-0000-000000000005', '50000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', 'Sopir cancel order', 1);

-- Insert data untuk tabel payments
INSERT INTO payments (payment_id, order_id, admin_fee, discount, status) VALUES
('80000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 2000.000, 5000.000, 'paid'),
('80000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', 2000.000, 0, 'pending'),
('80000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000003', 2000.000, 0, 'failed'),
('80000000-0000-0000-0000-000000000004', '50000000-0000-0000-0000-000000000004', 2000.000, 10000.000, 'paid'),
('80000000-0000-0000-0000-000000000005', '50000000-0000-0000-0000-000000000005', 2000.000, 0, 'pending');