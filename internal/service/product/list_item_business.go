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
	ListItem(ctx context.Context, key string) ([]productmodel.Product, error)
}

type listItemBusiness struct {
	store ListItemStorage
	cache ListItemCache
}

func NewListItemBiz(store ListItemStorage, cache ListItemCache) *listItemBusiness {
	return &listItemBusiness{store, cache}
}

func (biz *listItemBusiness) ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error) {
	// Tạo bản sao của Paging và Filter để sử dụng cho việc tạo khóa cache
	pagingCopy := *pagging
	filterCopy := *filter

	// Tạo khóa cache
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: productmodel.EntityName,
		Cond:       cond,
		Paging:     pagingCopy,
		Filter:     filterCopy,
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(productmodel.EntityName, err)
	}

	// Gọi cache với khóa đã tạo
	records, err := biz.cache.ListItem(ctx, key)
	if err != nil {
		return nil, common.ErrCannotListEntity(productmodel.EntityName, err)
	}

	if len(records) != 0 {
		return records, nil
	}

	// Gọi store để lấy dữ liệu từ database
	records, err = biz.store.ListItem(ctx, cond, pagging, filter, morekeys...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFoundEntity(productmodel.EntityName, err)
		}
		return nil, common.ErrCannotListEntity(productmodel.EntityName, err)
	}

	// Lưu vào cache với cùng khóa
	if len(records) != 0 {
		if err := biz.cache.SaveProduct(ctx, records, key); err != nil {
			return nil, common.ErrCannotCreateEntity(productmodel.EntityName, err)
		}
	}

	return records, nil
}
