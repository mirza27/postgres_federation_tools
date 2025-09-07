TRUNCATE TABLE
    "users",
    "drivers",
    "customers",
    "orders",
    "assignments",
    "reviews",
    "payments"
RESTART IDENTITY CASCADE;

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "drivers" CASCADE;
DROP TABLE IF EXISTS "customers" CASCADE;
DROP TABLE IF EXISTS "orders" CASCADE;
DROP TABLE IF EXISTS "assignments" CASCADE;
DROP TABLE IF EXISTS "reviews" CASCADE;
DROP TABLE IF EXISTS "payments" CASCADE;