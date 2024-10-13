CREATE TABLE Supplier (
                          supplier_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                          name VARCHAR(200) NOT NULL,
                          description TEXT,
                          contactInfo VARCHAR(200) NOT NULL,
                          address VARCHAR(200) NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                          created_by CHAR(30) DEFAULT 'system',
                          updated_by CHAR(30),
                          deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE SupplyOrder (
                             supplyorder_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                             order_date DATETIME NOT NULL,
                             description TEXT,
                             total_amount DECIMAL NOT NULL,
                             supplier_id BIGINT NOT NULL,
                             created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                             created_by CHAR(30) DEFAULT 'system',
                             updated_by CHAR(30),
                             deleted_at TIMESTAMP NULL DEFAULT NULL,
                             FOREIGN KEY (supplier_id) REFERENCES Supplier(supplier_id)
);

CREATE TABLE SupplyOrderItem (
                                 supplyorderitem_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                                 price DECIMAL NOT NULL,
                                 quantity INT NOT NULL,
                                 unit VARCHAR(200) NOT NULL,
                                 ingredient_id BIGINT NOT NULL,
                                 supplyorder_id BIGINT NOT NULL,
                                 stockbatch_id BIGINT NOT NULL,
                                 created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                 status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                                 created_by CHAR(30) DEFAULT 'system',
                                 updated_by CHAR(30),
                                 deleted_at TIMESTAMP NULL DEFAULT NULL,
                                 FOREIGN KEY (ingredient_id) REFERENCES Ingredient(ingredient_id),
                                 FOREIGN KEY (supplyorder_id) REFERENCES SupplyOrder(supplyorder_id),
                                 FOREIGN KEY (stockbatch_id) REFERENCES StockBatch(stockbatch_id)
);