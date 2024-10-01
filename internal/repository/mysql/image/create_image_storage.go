package imagestorage

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonrecover "tart-shop-manager/internal/common/recover"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	responseutil "tart-shop-manager/internal/util/response"
)

func (s *mysqlImage) CreateImage(ctx context.Context, data *imagemodel.CreateImage, morekeys ...string) (uint64, error) {
	db := s.db.Begin()
	if db.Error != nil {
		return 0, common.ErrDB(db.Error)
	}
	defer commonrecover.RecoverTransaction(db)

	// Step 1: Check if the URL already exists
	var existingImage imagemodel.Image
	err := db.WithContext(ctx).
		Where("url = ?", data.URL).
		First(&existingImage).Error

	if err == nil {
		// URL exists, return existing ImageID
		db.Rollback()
		return existingImage.ImageID, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// An unexpected error occurred
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	// Step 2: URL does not exist, create a new record
	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			// Handle duplicate entry error just in case
			fieldName := responseutil.ExtractFieldFromError(err, imagemodel.EntityName)
			db.Rollback()
			return 0, common.ErrDuplicateEntry(imagemodel.EntityName, fieldName, err)
		}
		db.Rollback()
		return 0, err
	}

	// Commit the transaction
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	return data.ImageID, nil
}
