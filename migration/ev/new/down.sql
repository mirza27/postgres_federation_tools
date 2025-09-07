TRUNCATE TABLE
    "drivers",
    "vehicles",
    "driver_vehicle_assignments",
    "shipments",
    "shipment_status_history",
    "shipment_packages"
RESTART IDENTITY CASCADE;

DROP TABLE IF EXISTS "drivers" CASCADE;
DROP TABLE IF EXISTS "vehicles" CASCADE;
DROP TABLE IF EXISTS "driver_vehicle_assignments" CASCADE;
DROP TABLE IF EXISTS "shipments" CASCADE;
DROP TABLE IF EXISTS "shipment_status_history" CASCADE;
DROP TABLE IF EXISTS "shipment_packages" CASCADE;