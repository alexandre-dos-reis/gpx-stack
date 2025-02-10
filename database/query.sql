-- name: FindAllProducts :many
SELECT
    *
FROM
    products;

-- name: FindOneProductBySlug :one
SELECT
    *
FROM
    products
WHERE
    slug = $1;

