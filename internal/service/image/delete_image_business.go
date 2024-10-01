package imagebusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
)

type DeleteImageStorage interface {
	DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type deleteImageBusiness struct {
	store DeleteImageStorage
}

func NewDeleteImageBiz(store DeleteImageStorage) *deleteImageBusiness {
	return &deleteImageBusiness{store: store}
}

func (biz *deleteImageBusiness) DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	if err := biz.store.DeleteImage(ctx, cond, morekeys...); err != nil {
		return common.ErrCannotDeleteEntity(imagemodel.EntityName, err)
	}

	return nil
}
