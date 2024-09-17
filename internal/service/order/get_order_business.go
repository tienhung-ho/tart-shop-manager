package orderbusiness

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetOrderStorage interface {
	GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error)
}

type GetOrderCache interface {
	GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error)
	SaveOrder(ctx context.Context, data interface{}, morekeys ...string) error
}

type getOrderBusiness struct {
	store GetOrderStorage
	cache GetOrderCache
}

func NewGetOrderBiz(store GetOrderStorage, cache GetOrderCache) *getOrderBusiness {
	return &getOrderBusiness{store, cache}
}

func (biz *getOrderBusiness) GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error) {
	record, err := biz.cache.GetOrder(ctx, cond, morekeys...)

	if err != nil {
		return nil, common.ErrNotFoundEntity(ordermodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetOrder(ctx, cond, morekeys...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil,
				common.ErrNotFoundEntity(ordermodel.EntityName, err)
		}

		return nil, common.ErrCannotGetEntity(ordermodel.EntityName, err)
	}

	if record != nil {

		var pagging paggingcommon.Paging
		pagging.Process()

		key := cacheutil.GenerateKey(ordermodel.EntityName, cond, pagging, commonfilter.Filter{})

		createOrder := record.ToCreateOrder()
		if err := biz.cache.SaveOrder(ctx, createOrder, key); err != nil {
			return nil, common.ErrCannotCreateEntity(ordermodel.EntityName, err)
		}
	}

	return record, err
}
