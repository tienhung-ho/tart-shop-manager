CREATE TABLE RecipeIngredient (
                                  recipe_id BIGINT NOT NULL,
                                  ingredient_id BIGINT NOT NULL,
                                  quantity DECIMAL(10,2) NOT NULL, -- Số lượng nguyên liệu, ví dụ: 100 (gram)

--     -- Common Fields
--                                   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--                                   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--                                   status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
--                                   created_by CHAR(30) DEFAULT 'system',
--                                   updated_by CHAR(30),
--                                   deleted_at TIMESTAMP NULL DEFAULT NULL,

                                  PRIMARY KEY (recipe_id, ingredient_id),
                                  FOREIGN KEY (recipe_id) REFERENCES Recipe(recipe_id),
                                  FOREIGN KEY (ingredient_id) REFERENCES Ingredient(ingredient_id)
--                                   INDEX idx_deleted_at (deleted_at) -- Add index for soft deletion
);
