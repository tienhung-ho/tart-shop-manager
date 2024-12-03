ALTER TABLE Product DROP COLUMN quantity_in_stock;
ALTER TABLE RecipeIngredient
    ADD COLUMN unit ENUM('kg', 'l', 'p') NOT NULL DEFAULT 'kg';
