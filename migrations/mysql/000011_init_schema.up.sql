ALTER TABLE order_product
    RENAME TO product_order;

-- Thêm các index cho bảng product_order sau khi đổi tên
ALTER TABLE
    ADD INDEX idx_order_id (order_id),
    ADD INDEX idx_product_id (product_id),
    ADD INDEX idx_price (price),
    ADD INDEX idx_order_product_price (order_id, product_id, price);

-- Thêm index cho bảng Order
ALTER TABLE `Order`
    ADD INDEX idx_status (status),
    ADD INDEX idx_created_at (created_at),
    ADD INDEX idx_account_status (account_id, status),
    ADD INDEX idx_total_amount (total_amount);

-- Thêm index cho bảng Product
ALTER TABLE Product
    ADD INDEX idx_status (status),
    ADD INDEX idx_category_status (category_id, status),
    ADD INDEX idx_name (name),
    ADD INDEX idx_quantity (quantity_in_stock);

ALTER TABLE Product
    ADD COLUMN price DECIMAL(11, 2) NOT NULL DEFAULT 0.00
    AFTER description;

-- 6_drop_order_product_table.sql
DROP TABLE IF EXISTS `order_product`;


-- 4_create_order_recipe_table.sql
CREATE TABLE `OrderRecipe` (
                                order_id BIGINT NOT NULL,
                                recipe_id BIGINT NOT NULL,
                                quantity INT NOT NULL DEFAULT 1,      -- Số lượng của recipe trong đơn hàng
                                price DECIMAL(11, 2) NOT NULL,        -- Giá recipe tại thời điểm đặt hàng
                                PRIMARY KEY (order_id, recipe_id),
                                FOREIGN KEY (order_id) REFERENCES `Order`(order_id)
                                    ON UPDATE CASCADE
                                    ON DELETE CASCADE,
                                FOREIGN KEY (recipe_id) REFERENCES `Recipe`(recipe_id)
                                    ON UPDATE CASCADE
                                    ON DELETE RESTRICT
);

ALTER TABLE `StockBatch`
    MODIFY COLUMN quantity FLOAT;

ALTER TABLE `Recipe`
    ADD CONSTRAINT unique_product_size UNIQUE (product_id, size);


