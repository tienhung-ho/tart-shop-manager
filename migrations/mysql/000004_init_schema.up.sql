-- 1_create_order_table.sql
CREATE TABLE `Order` (
                         order_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                         account_id BIGINT UNSIGNED NOT NULL,
                         total_amount DECIMAL(11, 2) NOT NULL DEFAULT 0.00,
                         tax DECIMAL(10, 2) DEFAULT 0.00,
                         status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                         created_by CHAR(30) DEFAULT NULL,
                         updated_by CHAR(30) DEFAULT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         deleted_at DATETIME DEFAULT NULL,
                         FOREIGN KEY (account_id) REFERENCES Account(account_id)
);

-- 2_create_order_product_table.sql
CREATE TABLE `order_product` (
                                 order_id BIGINT NOT NULL,
                                 product_id BIGINT NOT NULL,
                                 quantity INT NOT NULL DEFAULT 1,  -- Số lượng của sản phẩm trong đơn hàng
                                 price DECIMAL(11, 2) NOT NULL,   -- Giá sản phẩm tại thời điểm đặt hàng
                                 PRIMARY KEY (order_id, product_id),
                                 FOREIGN KEY (order_id) REFERENCES `Order`(order_id) ON UPDATE CASCADE,
                                 FOREIGN KEY (product_id) REFERENCES Product(product_id) ON UPDATE CASCADE
);

-- 3_remove_price_from_recipe.sql
ALTER TABLE Recipe DROP COLUMN price;
