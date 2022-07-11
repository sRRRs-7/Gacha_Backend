CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar UNIQUE NOT NULL,
  "hash_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar UNIQUE NOT NULL,
  "balance" bigint NOT NULL DEFAULT 100,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "gachas" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "items" (
  "id" bigserial PRIMARY KEY,
  "item_name" varchar NOT NULL,
  "rating" int NOT NULL,
  "item_url" varchar UNIQUE NOT NULL,
  "category_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "category" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "galleries" (
  "id" bigserial PRIMARY KEY,
  "owner_id" bigint NOT NULL,
  "item_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "sessions" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expired_at" timestamptz NOT NULL
);

CREATE INDEX ON "galleries" USING BTREE ("owner_id");

CREATE INDEX ON "galleries" USING BTREE ("item_id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("user_name");

ALTER TABLE "gachas" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "gachas" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "items" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "galleries" ADD FOREIGN KEY ("owner_id") REFERENCES "accounts" ("id");

ALTER TABLE "galleries" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_name") REFERENCES "users" ("user_name");
