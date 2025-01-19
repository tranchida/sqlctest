CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL
);

ALTER TABLE product ADD column color VARCHAR(255);

CREATE TABLE customer (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE "order" (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customer(id),
    FOREIGN KEY (product_id) REFERENCES product(id)
);

INSERT INTO customer (name,email) VALUES ('John Doe','jdoe@fake.net');

INSERT INTO product (code,price,stock,color) VALUES ('sample1',30,81,'white');
INSERT INTO product (code,price,stock,color) VALUES ('sample2',20,72,'red');
INSERT INTO product (code,price,stock,color) VALUES ('sample3',30,43,'blue');
INSERT INTO product (code,price,stock,color) VALUES ('sample4',40,37,'green');
INSERT INTO product (code,price,stock,color) VALUES ('sample5',50,91,'yellow');

INSERT INTO "order" (customer_id,product_id,quantity) VALUES (1,1,1);
INSERT INTO "order" (customer_id,product_id,quantity) VALUES (1,2,2);

