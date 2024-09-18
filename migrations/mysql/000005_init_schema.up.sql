CREATE TABLE Ingredient (
                            ingredient_id BIGINT NOT NULL AUTO_INCREMENT,
                            name VARCHAR(200) NOT NULL,
                            description TEXT,
                            unit VARCHAR(100) NOT NULL,

    -- Common Fields
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                            created_by CHAR(30) DEFAULT 'system',
                            updated_by CHAR(30),
                            deleted_at TIMESTAMP NULL DEFAULT NULL,

                            PRIMARY KEY (ingredient_id),
                            INDEX idx_deleted_at (deleted_at) -- Add index for soft deletion
);


CREATE TABLE StockBatch (
                            stockbatch_id BIGINT NOT NULL AUTO_INCREMENT,
                            quantity INT NOT NULL,
                            expiration_date DATE NOT NULL,
                            received_date DATE NOT NULL,
                            ingredient_id BIGINT NOT NULL,

    -- Common Fields
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                            created_by CHAR(30) DEFAULT 'system',
                            updated_by CHAR(30),
                            deleted_at TIMESTAMP NULL DEFAULT NULL,

                            PRIMARY KEY (stockbatch_id),
                            FOREIGN KEY (ingredient_id) REFERENCES Ingredient(ingredient_id) ON DELETE CASCADE,
                            INDEX idx_deleted_at (deleted_at) -- Add index for soft deletion
);
