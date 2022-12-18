CREATE TABLE
    "session" (
        "id" uuid PRIMARY KEY,
        "user_id" serial8 NOT NULL,
        "refresh_token" varchar(128) NOT NULL,
        "expires_at" timestamptz NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE
    "session_ip" (
        "id" serial8 PRIMARY KEY,
        "session_id" uuid NOT NULL,
        "user_agent" varchar(64),
        "ip" varchar(32) NOT NULL
    );

ALTER TABLE "session_ip"
ADD
    FOREIGN KEY ("session_id") REFERENCES "session" ("id");