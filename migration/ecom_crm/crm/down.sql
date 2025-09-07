-- delete data and its seq
TRUNCATE TABLE
    "tickets" ,
    "customers" ,
    "ticket_categories" ,
    "agents" 
RESTART IDENTITY CASCADE;

-- drop table
DROP TABLE IF EXISTS tickets CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS ticket_categories CASCADE;
DROP TABLE IF EXISTS agents CASCADE;