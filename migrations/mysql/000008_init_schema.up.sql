ALTER TABLE Recipe ADD COLUMN cost DECIMAL(10,2) NOT NULL;
ALTER TABLE Recipe ADD COLUMN status ENUM('Pending', 'Active', 'Inactive') DEFAULT 'Pending';