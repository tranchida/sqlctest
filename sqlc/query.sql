-- name: CountProducts :one
SELECT COUNT(*) FROM product;

-- name: GetProductById :one
SELECT * FROM product WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM product ORDER BY code;

-- name: CreateProduct :one
INSERT INTO product (code, price, stock) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateProduct :exec
UPDATE product SET code = $1, price = $2, stock = $3 WHERE id = $4 RETURNING *;

-- name: UpdateProductPrice :exec
UPDATE product SET price = $1 WHERE id = $2 RETURNING *;

-- name: CreateCustomer :one
INSERT INTO customer (name, email) VALUES ($1, $2) RETURNING *;

-- name: CreateOrder :one
INSERT INTO "order" (customer_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING *;

-- name: ListOrders :many
SELECT 
    customer.*, 
    product.*, 
    "order".quantity 
FROM customer 
join "order" on customer.id = "order".customer_id
join product on product.id = "order".product_id
ORDER BY customer.name, product.code;