-- Pengemudi
INSERT INTO pengemudi (nama_lengkap, no_telepon, tgl_lahir, alamat_tinggal) VALUES
('Budi Santoso', '+628112345678', '1982-05-19', 'Jl. Merdeka No. 12, Jakarta'),
('Dewi Lestari', '+628112345679', '1991-12-03', 'Jl. Sudirman No. 45, Bandung'),
('Irfan Jaya', '+628112345680', '1987-07-14', 'Jl. Gatot Subroto No. 8, Surabaya');

-- Kendaraan (One-to-Many)
INSERT INTO kendaraan (id_pengemudi, no_plat, jenis_kendaraan, merk, tahun_pembuatan, kapasitas_muat) VALUES
(1, 'B 1234 ABC', 'Pickup', 'Toyota', 2018, 1000.00),
(2, 'D 5678 XYZ', 'Truck', 'Mitsubishi', 2015, 4000.00),
(3, 'L 9012 QRS', 'Minivan', 'Honda', 2020, 800.00);

-- Pengiriman
INSERT INTO pengiriman (id_pengemudi, id_kendaraan, tgl_pengiriman, alamat_tujuan, status_pengiriman) VALUES
(1, 1, '2023-10-10', 'Jl. Thamrin No. 21, Jakarta', 'SELESAI'),
(2, 2, '2023-10-11', 'Jl. Asia Afrika No. 67, Bandung', 'DIKIRIM'),
(3, 3, '2023-10-12', 'Jl. Pemuda No. 33, Surabaya', 'DRAFT');

-- Status Pengiriman
INSERT INTO status_pengiriman (id_pengiriman, status, catatan) VALUES
(1, 'DRAFT', 'Menunggu konfirmasi'),
(1, 'DIKIRIM', 'Paket telah diambil'),
(1, 'SELESAI', 'Terkirim ke penerima'),
(2, 'DRAFT', 'Sedang diproses gudang');