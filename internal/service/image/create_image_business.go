package imagebusiness

import (
	"context"
	"log"
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	cloudutil "tart-shop-manager/internal/util/cloud"
)

type CreateImageStorage interface {
	CreateImage(ctx context.Context, data *imagemodel.CreateImage, morekeys ...string) (uint64, error)
}

type createImageBusiness struct {
	store CreateImageStorage
}

func NewCreateImageBiz(store CreateImageStorage) *createImageBusiness {
	return &createImageBusiness{store: store}
}

func (biz *createImageBusiness) CreateImage(ctx context.Context, data *cloudutil.Image, morekeys ...string) (uint64, string, error) {

	// Bước 1: Upload ảnh lên cloud
	fileURL, err := cloudutil.UploadSingleImageToS3(ctx, data.FileBuffer, data.FileName)
	if err != nil {
		return 0, "", common.ErrCannotUploadFile("image", err)
	}

	// Bước 2: Lưu thông tin ảnh vào cơ sở dữ liệu
	createImage := &imagemodel.CreateImage{
		URL:     fileURL,
		AltText: data.FileName,
	}

	recordID, err := biz.store.CreateImage(ctx, createImage, fileURL)
	if err != nil {
		// Nếu lưu vào DB thất bại, cần xóa ảnh đã upload để duy trì tính nhất quán
		deleteErr := cloudutil.DeleteSingleImageFromS3(ctx, data.FileName)
		if deleteErr != nil {
			log.Println(deleteErr.Error())
		}
		return 0, "", common.ErrCannotCreateEntity(imagemodel.EntityName, err)
	}

	return recordID, fileURL, nil
}
