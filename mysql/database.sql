CREATE TABLE products (
    product_id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    author VARCHAR(50),
    collection VARCHAR(50),
    image_url  VARCHAR(100),
    price DOUBLE NOT NULL,
    sales DOUBLE ,
    publisher VARCHAR(50),
    date_publish DATE ,
    rate DOUBLE ,
    description TEXT,
    status_product BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);
CREATE TABLE in_stocks (
    product_id BIGINT UNSIGNED,
    amount_sold BIGINT NOT NULL,
    amount_stock BIGINT NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);

CREATE TABLE inventory (
    inventory_id BIGINT UNSIGNED PRIMARY KEY,
    product_id BIGINT UNSIGNED NOT NULL,
    purchase_price DOUBLE NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id)
    REFERENCES products(product_id)
);	

CREATE TABLE users (
    user_id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_name VARCHAR(20) NOT NULL ,
    password VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL,
    role VARCHAR(15) DEFAULT 'CUSTOMER',
    first_name VARCHAR(20) ,
    last_name VARCHAR(20) ,
    address VARCHAR(50),
    birthday DATE,
    sex CHAR(10),
    phone CHAR(15),
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
    order_id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED,
    total_price DOUBLE NOT NULL,
    shipping_method VARCHAR(20) NOT NULL,
    receiver_address TEXT NOT NULL,
    receiver_phone CHAR(10),
    note TEXT,
    pay_method VARCHAR(20),
    number_cart VARCHAR(20),
    owner_name	VARCHAR(50),
    shipped_date DATE,
    status VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE order_items (
    order_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    amount INT UNSIGNED NOT NULL,
    unit_price DOUBLE,
    FOREIGN KEY (order_id) REFERENCES orders(order_id)
);

CREATE TABLE review (
    product_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    comment TEXT,
    rate INT UNSIGNED NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(user_id),
    FOREIGN KEY (product_id)
    REFERENCES products(product_id)

);
