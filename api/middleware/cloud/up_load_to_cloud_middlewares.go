// internal/cloudmiddleware/image_validation_middleware.go

package cloudmiddleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"

	"github.com/gin-gonic/gin"
)

func ImageValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			c.Abort()
			return
		}

		// Giới hạn kích thước file (ví dụ: 100MB)
		if file.Size > 100*1024*1024 {
			c.JSON(http.StatusBadRequest, common.ErrFileTooLarge(imagemodel.EntityName, fmt.Errorf("file size exceeds 100MB")))
			c.Abort()
			return
		}

		// Kiểm tra định dạng file (ví dụ: chỉ cho phép jpg, png)
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
		}
		if !allowedTypes[file.Header.Get("Content-Type")] {
			c.JSON(http.StatusBadRequest, common.ErrUnsupportedFileType(imagemodel.EntityName, fmt.Errorf("unsupported file type")))
			c.Abort()
			return
		}

		// Đọc file vào buffer và lưu vào context để handler sử dụng
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
			c.Abort()
			return
		}
		defer src.Close()

		fileBuffer := new(bytes.Buffer)
		_, err = io.Copy(fileBuffer, src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrInternal(err))
			c.Abort()
			return
		}

		c.Set("fileBuffer", fileBuffer.Bytes())
		c.Set("fileName", file.Filename)

		c.Next()
	}
}
