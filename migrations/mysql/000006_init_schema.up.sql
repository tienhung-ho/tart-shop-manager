CREATE TABLE Image (
    image_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    url VARCHAR(300) NOT NULL UNIQUE ,
    alt_text VARCHAR(255),
);


ALTER TABLE Category
    ADD COLUMN image_id BIGINT DEFAULT NULL,  -- Thêm cột image_id để lưu id của ảnh
ADD CONSTRAINT fk_category_image -- Tạo khóa ngoại
FOREIGN KEY (image_id) REFERENCES Image(image_id); -- Khóa ngoại tham chiếu đến image_id trong bảng Image


ALTER TABLE Product
    ADD COLUMN image_id BIGINT DEFAULT NULL,  -- Thêm cột image_id để tham chiếu đến bảng Image
DROP COLUMN image_url,  -- Xóa cột image_url nếu không còn cần thiết
ADD CONSTRAINT fk_product_image -- Tạo khóa ngoại
FOREIGN KEY (image_id) REFERENCES Image(image_id); -- Khóa ngoại tham chiếu đến image_id trong bảng Image
