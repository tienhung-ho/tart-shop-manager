package categorybusiness

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
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
		// Generate cache key
		key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
			EntityName: categorymodel.EntityName,
			Cond:       cond,
			Paging:     pagging,
			Filter:     commonfilter.Filter{},
			MoreKeys:   morekeys,
		})
		if err != nil {
			return nil, common.ErrCannotGenerateKey(categorymodel.EntityName, err)
		}

		log.Print("123213213", key)

		if err := biz.cache.SaveCategory(ctx, createCategore, key); err != nil {
			return nil, common.ErrCannotCreateEntity(categorymodel.EntityName, err)
		}

	}

	return record, nil
}
