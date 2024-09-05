CREATE TABLE `Role` (
                        role_id INT AUTO_INCREMENT PRIMARY KEY,
                        name VARCHAR(255) NOT NULL UNIQUE,
                        description TEXT,
                        status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                        `created_by` CHAR(30) DEFAULT NULL,
                        `updated_by` CHAR(30) DEFAULT NULL,
                        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        `deleted_at` DATETIME DEFAULT NULL   -- Thêm trường deleted_at cho Soft Delete
);

CREATE TABLE `Permission` (
                              permission_id INT AUTO_INCREMENT PRIMARY KEY,
                              name VARCHAR(255) NOT NULL UNIQUE,
                              description TEXT,
                              status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
                              `created_by` CHAR(30) DEFAULT NULL,
                              `updated_by` CHAR(30) DEFAULT NULL,
                              `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                              `deleted_at` DATETIME DEFAULT NULL   -- Thêm trường deleted_at cho Soft Delete
);

CREATE TABLE `role_permissions` (
                                    role_id INT,
                                    permission_id INT,
                                    PRIMARY KEY (role_id, permission_id),
                                    FOREIGN KEY (role_id) REFERENCES `Role`(role_id),   -- Tham chiếu đúng khoá chính role_id
                                    FOREIGN KEY (permission_id) REFERENCES `Permission`(permission_id)   -- Tham chiếu đúng khoá chính permission_id
);
