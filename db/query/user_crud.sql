-- name: CreateUser :exec
INSERT INTO "user" (
    "email", "user_name", "hashed_password"
) VALUES (
    $1, $2, $3
);

-- name: GetAllUser :many
SELECT * FROM "user";