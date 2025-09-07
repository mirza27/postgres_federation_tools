-- delete data and its seq
TRUNCATE TABLE 
    "payments",
    "order_items",
    "orders",
    "products",
    "customers",
    "product_categories"
RESTART IDENTITY CASCADE;

-- drop table
DROP TABLE IF EXISTS "payments" CASCADE;
DROP TABLE IF EXISTS "order_items" CASCADE;
DROP TABLE IF EXISTS "orders" CASCADE;
DROP TABLE IF EXISTS "products" CASCADE;
DROP TABLE IF EXISTS "customers" CASCADE;
DROP TABLE IF EXISTS "product_categories" CASCADE;