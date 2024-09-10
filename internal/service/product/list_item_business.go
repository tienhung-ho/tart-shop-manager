package productbusiness

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type ListItemStorage interface {
	ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error)
}

type ListItemCache interface {
	SaveProduct(ctx context.Context, data interface{}, morekeys ...string) error
	ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error)
}

type listItemBusiness struct {
	store ListItemStorage
	cache ListItemCache
}

func NewListItemBiz(store ListItemStorage, cache ListItemCache) *listItemBusiness {
	return &listItemBusiness{store, cache}
}

func (biz *listItemBusiness) ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error) {

	records, err := biz.cache.ListItem(ctx, cond, pagging, filter, morekeys...)

	if err != nil {
		return nil, common.ErrCannotListEntity(productmodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	records, err = biz.store.ListItem(ctx, cond, pagging, filter, morekeys...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, common.ErrNotFoundEntity(productmodel.EntityName, err)
		}

		return nil, common.ErrCannotListEntity(productmodel.EntityName, err)
	}

	if len(records) != 0 {

		key := cacheutil.GenerateKey(productmodel.EntityName, cond, *pagging, *filter)
		err := biz.cache.SaveProduct(ctx, records, key)

		if err != nil {
			return nil, common.ErrCannotCreateEntity(productmodel.EntityName, err)
		}
	}

	return records, nil
}
