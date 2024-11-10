ALTER TABLE `Order`
    ADD COLUMN `order_date` DATETIME NULL AFTER `account_id`;


UPDATE `Order` SET `order_date` = `created_at` WHERE `order_date` IS NULL;

ALTER TABLE `Order`
    MODIFY COLUMN `order_date` DATETIME NOT NULL;


ALTER TABLE `StockBatch`
    MODIFY `expiration_date` DATETIME NULL,
    MODIFY `received_date` DATETIME NULL;


ALTER TABLE Product DROP FOREIGN KEY fk_product_image;

-- Xóa cột image_id
ALTER TABLE Product DROP COLUMN image_id;