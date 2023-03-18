-- name: CreateUser :exec
INSERT INTO "user" (
    "email", "user_name", "hashed_password"
) VALUES (
    $1, $2, $3
);

-- name: GetAllUser :many
SELECT * FROM "user";

-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1  LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM "user"
WHERE "user".id = $1  LIMIT 1;

-- name: RegisterSupplier :exec
UPDATE "user" SET "role" = 'supplier' WHERE "id" = $1;

-- name: GetPhone :one
SELECT * FROM "user_profile"
WHERE "user_id" = $1;

-- name: GetAllAddress :many
SELECT * FROM "user_address"
WHERE "user_id" = $1;

-- name: UpdateUserName :exec
UPDATE "user"
SET "user_name" = $1,
"phone" = $2,
"address" = $4,
"note" = $5
WHERE "id" = $3;

