package categorybusiness

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetCategoryStorage interface {
	GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
}

type GetCategoryCache interface {
	GetCategory(ctx context.Context,
		cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error)
	SaveCategory(ctx context.Context, data interface{}, morekeys ...string) error
}

type getCategoryBusiness struct {
	store GetCategoryStorage
	cache GetCategoryCache
}

func NewGetCategoryBiz(store GetCategoryStorage, cache GetCategoryCache) *getCategoryBusiness {
	return &getCategoryBusiness{store, cache}
}

func (biz *getCategoryBusiness) GetCategory(ctx context.Context,
	cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error) {

	record, err := biz.cache.GetCategory(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrNotFoundEntity(categorymodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetCategory(ctx, cond, morekeys...)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil,
				common.ErrNotFoundEntity(categorymodel.EntityName, err)
		}

		return nil, common.ErrCannotGetEntity(categorymodel.EntityName, err)
	}

	if record != nil {

		var pagging paggingcommon.Paging
		pagging.Process()

		createCategore := record.ToCreateCategoryCache()
		key := cacheutil.GenerateKey(categorymodel.EntityName, cond, pagging, commonfilter.Filter{})
		if err := biz.cache.SaveCategory(ctx, createCategore, key); err != nil {
			return nil, common.ErrCannotCreateEntity(categorymodel.EntityName, err)
		}

	}

	return record, nil
}
