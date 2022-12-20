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

-- name: RegisterSupplier :exec
UPDATE "user" SET "role" = 'supplier' WHERE "id" = $1;