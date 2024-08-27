CREATE TABLE `Account` (
    `account_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `role_id` TINYINT UNSIGNED NOT NULL,
    `phone` VARCHAR(20) NOT NULL UNIQUE,
    `fullname` VARCHAR(300) DEFAULT NULL,
    `avatar_url` VARCHAR(255) DEFAULT NULL,
    `password` VARCHAR(200) NOT NULL,
    `email` VARCHAR(100) NOT NULL UNIQUE,
    `gender` ENUM('Male', 'Female', 'Other') DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `status` ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending',
    `created_by` CHAR(30) DEFAULT NULL,
    `updated_by` CHAR(30) DEFAULT NULL,
     PRIMARY KEY (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
