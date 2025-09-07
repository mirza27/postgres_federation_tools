-- Product Categories
INSERT INTO product_categories (name, description) VALUES
('Elektronik', 'Perangkat elektronik dan gadget terbaru'),
('Fashion', 'Pakaian dan aksesoris terkini'),
('Makanan', 'Makanan dan minuman'),
('Kesehatan', 'Produk kesehatan dan kecantikan'),
('Rumah Tangga', 'Perabotan dan perlengkapan rumah');

-- Products
INSERT INTO products (name, price, description, product_category_id) VALUES
('Smartphone X', 8500000, 'RAM 8GB, Kamera 48MP', 1),
('Laptop Pro', 12500000, 'Processor i7, SSD 512GB', 1),
('Kemeja Slim Fit', 450000, 'Katun premium, berbagai warna', 2),
('Kopi Arabika', 120000, 'Biji kopi premium 500g', 3),
('Blender Multifungsi', 650000, '5 Pisau stainless steel', 5),
('Vitamin C', 150000, 'Suplemen kesehatan 60 tablet', 4),
('Sofa Minimalis', 2500000, 'Bahan kulit sintetis', 5);

-- Customers
INSERT INTO customers (name, email, phone, address) VALUES
('Budi Santoso', 'budi@mail.com', '08123456789', 'Jl. Merdeka No.12, Jakarta'),
('Siti Rahayu', 'siti@mail.com', '08234567890', 'Jl. Sudirman No.45, Bandung'),
('Ahmad Fauzi', 'ahmad@mail.com', '08345678901', 'Jl. Gatot Subroto No.78, Surabaya'),
('Dewi Anggraini', 'dewi@mail.com', '08456789012', 'Jl. Thamrin No.34, Medan'),
('Rudi Hermawan', 'rudi@mail.com', '08567890123', 'Jl. Asia Afrika No.56, Bandung');

-- Orders
INSERT INTO orders (customer_id, status, delivery_type, total_amount, grand_total) VALUES
(1, 'completed', 'express', 8950000, 9000000),
(2, 'processing', 'regular', 5000000, 5050000),
(3, 'pending', 'regular', 1200000, 1200000),
(4, 'shipped', 'express', 650000, 700000),
(1, 'completed', 'regular', 450000, 450000);

-- Order Items
INSERT INTO order_items (product_id, quantity, price_per_item, total_amount, "order_id") VALUES
(1, 1, 8500000, 8500000, 1),
(3, 1, 450000, 450000, 1),
(5, 1, 3500000, 3500000, 2),
(4, 5, 120000, 600000, 2),
(6, 2, 150000, 300000, 3),
(7, 1, 2500000, 2500000, 4),
(3, 1, 450000, 450000, 5);

-- Payments
INSERT INTO payments (customer_id, order_id, status, payment_fee, total_amount, grand_total) VALUES
(1, 1, 'success', 50000, 8950000, 9000000),
(2, 2, 'pending', 50000, 5000000, 5050000),
(3, 3, 'failed', 0, 1200000, 1200000),
(4, 4, 'success', 50000, 650000, 700000),
(1, 5, 'success', 0, 450000, 450000);