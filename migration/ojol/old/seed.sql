-- Insert data untuk tabel pelanggan (5 data)
INSERT INTO
    pelanggan (nama, telepon, email, password)
VALUES
    (
        'Ahmad Riyadi',
        '081234567890',
        'ahmad.riyadi@email.com',
        'hashed_password_1'
    ),
    (
        'Siti Nurhaliza',
        '081298765432',
        'siti.nurhaliza@email.com',
        'hashed_password_2'
    ),
    (
        'Budi Santoso',
        '081112223333',
        'budi.santoso@email.com',
        'hashed_password_3'
    ),
    (
        'Dewi Lestari',
        '081344445555',
        'dewi.lestari@email.com',
        'hashed_password_4'
    ),
    (
        'Rina Melati',
        '081566667777',
        'rina.melati@email.com',
        'hashed_password_5'
    );

-- Insert data untuk tabel ojol (5 data)
INSERT INTO
    ojol (nama, telepon, email, password, kendaraan, status)
VALUES
    (
        'Joko Widodo',
        '081677778888',
        'joko.widodo@email.com',
        'hashed_password_6',
        '{"type": "motor", "merk": "Honda", "plat": "B 1234 AB"}',
        'online'
    ),
    (
        'Surya Paloh',
        '081788889999',
        'surya.paloh@email.com',
        'hashed_password_7',
        '{"type": "motor", "merk": "Yamaha", "plat": "B 5678 CD"}',
        'offline'
    ),
    (
        'Bambang Pamungkas',
        '081899990000',
        'bambang.pamungkas@email.com',
        'hashed_password_8',
        '{"type": "mobil", "merk": "Toyota", "plat": "B 9012 EF"}',
        'online'
    ),
    (
        'Dian Sastro',
        '081900001111',
        'dian.sastro@email.com',
        'hashed_password_9',
        '{"type": "motor", "merk": "Suzuki", "plat": "B 3456 GH"}',
        'online'
    ),
    (
        'Taufik Hidayat',
        '081011112222',
        'taufik.hidayat@email.com',
        'hashed_password_10',
        '{"type": "mobil", "merk": "Daihatsu", "plat": "B 7890 IJ"}',
        'offline'
    );

-- Insert data untuk tabel pesanan (10 data)
INSERT INTO
    pesanan (
        id_pelanggan,
        id_ojol,
        status,
        nominal,
        catatan,
        titik_jemput,
        titik_tujuan
    )
VALUES
    (
        1,
        1,
        'selesai',
        25000.50,
        'Tolong antar cepat',
        'Jl. Merdeka No. 10',
        'Plaza Indonesia'
    ),
    (
        2,
        3,
        'dalam_perjalanan',
        45000.00,
        'Jangan ngebut',
        'Jl. Sudirman No. 25',
        'GBK Stadium'
    ),
    (
        3,
        2,
        'dibatalkan',
        30000.75,
        'Sudah ditunggu',
        'Jl. Thamrin No. 8',
        'Monas'
    ),
    (
        4,
        4,
        'selesai',
        60000.00,
        'Terima kasih',
        'Jl. Gatot Subroto No. 12',
        'Senayan City'
    ),
    (
        5,
        5,
        'menunggu',
        35000.00,
        'Tolong hati-hati',
        'Jl. Asia Afrika No. 100',
        'Braga Street'
    ),
    (
        1,
        3,
        'dalam_perjalanan',
        28000.50,
        'Pesanan penting',
        'Jl. Mangga Dua No. 5',
        'Ancol Dreamland'
    ),
    (
        2,
        1,
        'selesai',
        52000.00,
        'Terima kasih',
        'Jl. Kebon Jeruk No. 15',
        'Central Park'
    ),
    (
        3,
        4,
        'dibatalkan',
        42000.75,
        'Mohon cepat',
        'Jl. Palmerah No. 20',
        'Summarecon Mall'
    ),
    (
        4,
        2,
        'menunggu',
        38000.00,
        'Hati-hati di jalan',
        'Jl. Pasar Minggu No. 30',
        'Pondok Indah Mall'
    ),
    (
        5,
        5,
        'selesai',
        47000.50,
        'Terima kasih',
        'Jl. Casablanca No. 40',
        'Kuningan City'
    );

-- Insert data untuk tabel registrasi_ojol (5 data)
INSERT INTO
    registrasi_ojol (
        id_ojol,
        nama_lengkap,
        nik,
        alamat,
        status_registrasi
    )
VALUES
    (
        1,
        'Joko Widodo',
        '1234567890123456',
        'Jl. Kenangan No. 1, Jakarta',
        'disetujui'
    ),
    (
        2,
        'Surya Paloh',
        '2345678901234567',
        'Jl. Kenangan No. 2, Jakarta',
        'ditolak'
    ),
    (
        3,
        'Bambang Pamungkas',
        '3456789012345678',
        'Jl. Kenangan No. 3, Jakarta',
        'disetujui'
    ),
    (
        4,
        'Dian Sastro',
        '4567890123456789',
        'Jl. Kenangan No. 4, Jakarta',
        'menunggu'
    ),
    (
        5,
        'Taufik Hidayat',
        '5678901234567890',
        'Jl. Kenangan No. 5, Jakarta',
        'disetujui'
    );