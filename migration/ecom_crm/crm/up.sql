CREATE TABLE "customers" (
  "id" SERIAL PRIMARY KEY,
  "contact_name" varchar NOT NULL,
  "company_name" varchar,
  "phone" varchar NOT NULL,
  "email" varchar,
  "address" text,
  "city" varchar,
  "customer_type" varchar,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "ticket_categories" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "tickets" (
  "id" SERIAL PRIMARY KEY,
  "customer_id" int,
  "subject" varchar NOT NULL,
  "description" varchar NOT NULL,
  "ticket_category_id" int,
  "status" varchar,
  "priority" varchar,
  "agent_id" int,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "agents" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar,
  "status" varchar,
  "created_at" timestamptz DEFAULT (now())
);


ALTER TABLE "tickets" ADD FOREIGN KEY ("ticket_category_id") REFERENCES "ticket_categories" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("agent_id") REFERENCES "agents" ("id");

ALTER TABLE "tickets" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");
