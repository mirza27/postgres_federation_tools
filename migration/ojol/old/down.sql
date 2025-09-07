TRUNCATE TABLE
    "pelanggan",
    "ojol",
    "pesanan",
    "registrasi_ojol"
RESTART IDENTITY CASCADE;

DROP TABLE IF EXISTS "pelanggan" CASCADE;
DROP TABLE IF EXISTS "ojol" CASCADE;
DROP TABLE IF EXISTS "pesanan" CASCADE;
DROP TABLE IF EXISTS "registrasi_ojol" CASCADE;