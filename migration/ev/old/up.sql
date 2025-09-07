-- Tabel utama untuk pengemudi
CREATE TABLE pengemudi (
    id_pengemudi SERIAL PRIMARY KEY,
    nama_lengkap VARCHAR(100) NOT NULL,
    no_telepon VARCHAR(20),
    tgl_lahir DATE,
    alamat_tinggal TEXT,
    status_aktif BOOLEAN DEFAULT TRUE,
    tgl_dibuat TIMESTAMP DEFAULT NOW()
);

-- Tabel untuk kendaraan
CREATE TABLE kendaraan (
    id_kendaraan SERIAL PRIMARY KEY,
    id_pengemudi INTEGER REFERENCES pengemudi(id_pengemudi),
    no_plat VARCHAR(15) NOT NULL,
    jenis_kendaraan VARCHAR(50),
    merk VARCHAR(50),
    tahun_pembuatan INTEGER,
    kapasitas_muat DECIMAL(10,2)
);  

-- Tabel untuk pengiriman
CREATE TABLE pengiriman (
    id_pengiriman SERIAL PRIMARY KEY,
    id_pengemudi INTEGER REFERENCES pengemudi(id_pengemudi),
    id_kendaraan INTEGER REFERENCES kendaraan(id_kendaraan),
    tgl_pengiriman DATE,
    alamat_tujuan TEXT,
    status_pengiriman VARCHAR(20) DEFAULT 'DRAFT',
    catatan TEXT
);

-- Tabel untuk riwayat status pengiriman
CREATE TABLE status_pengiriman (
    id_status SERIAL PRIMARY KEY,
    id_pengiriman INTEGER REFERENCES pengiriman(id_pengiriman),
    status VARCHAR(20),
    tgl_status TIMESTAMP DEFAULT NOW(),
    catatan TEXT
);