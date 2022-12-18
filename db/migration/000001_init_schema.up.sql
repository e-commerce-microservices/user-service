CREATE TYPE user_role AS ENUM ('customer', 'supplier', 'admin');

CREATE TABLE
    "user" (
        "id" serial8 PRIMARY KEY,
        "email" varchar(128) UNIQUE NOT NULL,
        "role" user_role NOT NULL DEFAULT 'customer',
        "active_status" boolean NOT NULL DEFAULT FALSE,
        "hashed_password" varchar(128) NOT NULL,
        "password_updated_at" timestamptz NOT NULL DEFAULT (now()),
        "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE
    "user_profile" (
        "id" serial8 PRIMARY KEY,
        "user_id" serial8 NOT NULL,
        "user_name" varchar(64) NOT NULL,
        "phone" varchar(32),
        "avatar" varchar(256),
        "created_at" timestamptz NOT NULL DEFAULT (now()),
        "updated_at" timestamptz NOT NULL DEFAULT (now())
    );

ALTER TABLE "user_profile"
ADD
    FOREIGN KEY ("user_id") REFERENCES "user" ("id");

CREATE INDEX ON "user_profile" ("user_id");

CREATE TABLE
    "user_address" (
        "id" serial8 PRIMARY KEY,
        "user_id" serial8 NOT NULL,
        "address" varchar(256) NOT NULL,
        "note" varchar(256),
        "created_at" timestamptz NOT NULL DEFAULT (now()),
        "updated_at" timestamptz NOT NULL DEFAULT (now())
    );

ALTER TABLE "user_address"
ADD
    FOREIGN KEY ("user_id") REFERENCES "user" ("id");

CREATE INDEX ON "user_address" ("user_id");