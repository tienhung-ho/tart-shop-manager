ALTER TABLE Image
    ADD COLUMN product_id BIGINT DEFAULT NULL,
ADD CONSTRAINT fk_product
    FOREIGN KEY (product_id) REFERENCES Product(product_id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

ALTER TABLE Image
ADD COLUMN category_id BIGINT DEFAULT NULL,
ADD CONSTRAINT fk_category
    FOREIGN KEY (category_id) REFERENCES Category(category_id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

ALTER TABLE Image
    ADD COLUMN account_id BIGINT UNSIGNED DEFAULT NULL,
ADD CONSTRAINT fk_account
    FOREIGN KEY (account_id) REFERENCES Account(account_id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;
