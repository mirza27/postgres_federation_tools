-- Agents
INSERT INTO agents (name, email, status) VALUES
('Rina Wijaya', 'rina@crm.com', 'active'),
('Dito Pratama', 'dito@crm.com', 'active'),
('Bambang Susanto', 'bambang@crm.com', 'active'),
('Lina Hartati', 'lina@crm.com', 'inactive'),
('Fajar Setiawan', 'fajar@crm.com', 'active');

-- Ticket Categories
INSERT INTO ticket_categories (name) VALUES
('Technical Support'),
('Billing'),
('Sales Inquiry'),
('Product Complaint'),
('Account Management');
    
-- Customers
INSERT INTO customers (contact_name, company_name, phone, email, address, city, customer_type) VALUES
('Budi Santoso', 'PT Maju Jaya', '021-123456', 'budi@majujaya.com', 'Jl. Industri No.1', 'Jakarta', 'corporate'),
('Siti Rahayu', NULL, '08123456789', 'siti@personal.com', 'Jl. Melati No.5', 'Bandung', 'individual'),
('Ahmad Fauzi', 'CV Kreatif Abadi', '022-987654', 'ahmad@kreatifabadi.co.id', 'Jl. Karya No.22', 'Bandung', 'corporate'),
('Dewi Anggraini', 'PT Global Solusi', '031-456789', 'dewi@globalsolusi.com', 'Jl. Teknologi No.8', 'Surabaya', 'corporate'),
('Rudi Hermawan', 'UD Sentosa', '024-567890', 'rudi@sentosa.com', 'Jl. Damai No.3', 'Semarang', 'corporate');

-- Tickets
INSERT INTO tickets (customer_id, subject, description, ticket_category_id, status, priority, agent_id) VALUES
(1, 'Gagal Login', 'Tidak bisa login ke dashboard', 1, 'open', 'high', 1),
(2, 'Pembayaran Tertunda', 'Konfirmasi pembayaran invoice #INV-001', 2, 'pending', 'medium', 2),
(3, 'Permintaan Demo Produk', 'Minta jadwal demo fitur premium', 3, 'closed', 'low', 3),
(4, 'Produk Rusak', 'Barang datang dalam kondisi rusak', 4, 'open', 'critical', 1),
(5, 'Upgrade Paket', 'Ingin upgrade ke paket enterprise', 5, 'resolved', 'high', 4),
(1, 'Integrasi API', 'Error saat integrasi dengan sistem HR', 1, 'open', 'high', 5);