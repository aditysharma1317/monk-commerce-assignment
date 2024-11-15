CREATE TABLE IF NOT EXISTS coupons (
    id uuid PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cart_wise_coupons (
    coupon_id uuid PRIMARY KEY,
    threshold DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(5, 2) NOT NULL,
    FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_wise_coupons (
    coupon_id uuid PRIMARY KEY,
    product_id VARCHAR(255) NOT NULL,
    discount DECIMAL(5, 2) NOT NULL,
    FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bx_gy_coupons (
    coupon_id uuid PRIMARY KEY,
    repetition_limit INT NOT NULL,
    FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bx_gy_buy_products (
    bx_gy_coupon_id uuid,
    product_id uuid NOT NULL,
    quantity INT NOT NULL,
    PRIMARY KEY (bx_gy_coupon_id, product_id),
    FOREIGN KEY (bx_gy_coupon_id) REFERENCES bx_gy_coupons(coupon_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bx_gy_get_products (
    bx_gy_coupon_id uuid,
    product_id uuid NOT NULL,
    quantity INT NOT NULL,
    PRIMARY KEY (bx_gy_coupon_id, product_id),
    FOREIGN KEY (bx_gy_coupon_id) REFERENCES bx_gy_coupons(coupon_id) ON DELETE CASCADE
);

