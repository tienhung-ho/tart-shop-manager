CREATE TABLE Image (
                       image_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                       url VARCHAR(300) NOT NULL UNIQUE,
                       alt_text VARCHAR(255)
);

ALTER TABLE Category
    ADD COLUMN image_id BIGINT DEFAULT NULL,
    ADD CONSTRAINT fk_category_image
    FOREIGN KEY (image_id) REFERENCES Image(image_id);

ALTER TABLE Product
    ADD COLUMN image_id BIGINT DEFAULT NULL,
DROP COLUMN image_url,
    ADD CONSTRAINT fk_product_image
    FOREIGN KEY (image_id) REFERENCES Image(image_id);