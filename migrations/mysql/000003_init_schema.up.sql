CREATE TABLE Category (
                          category_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                          name VARCHAR(200) NOT NULL,
                          description TEXT,
                          status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                          created_by CHAR(30) DEFAULT NULL,
                          updated_by CHAR(30) DEFAULT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          deleted_at DATETIME DEFAULT NULL
);
CREATE TABLE Product (
                         product_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                         name VARCHAR(200) NOT NULL,
                         description TEXT,
                         quantity_in_stock INT NOT NULL,
                         image_url VARCHAR(300) NOT NULL,
                         category_id BIGINT NOT NULL,  -- Khóa ngoại để liên kết với Category
                         status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                         created_by CHAR(30) DEFAULT NULL,
                         updated_by CHAR(30) DEFAULT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         deleted_at DATETIME DEFAULT NULL,
                         FOREIGN KEY (category_id) REFERENCES Category(category_id)
);
CREATE TABLE Recipe (
                        recipe_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
                        product_id BIGINT NOT NULL,  -- Khóa ngoại để liên kết với Product
                        size ENUM('Small', 'Medium', 'Large') NOT NULL,
                        price DECIMAL(10,2) NOT NULL,
                        description TEXT,
                        created_by CHAR(30) DEFAULT NULL,
                        updated_by CHAR(30) DEFAULT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        deleted_at DATETIME DEFAULT NULL,
                        FOREIGN KEY (product_id) REFERENCES Product(product_id)
);
