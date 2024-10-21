package supplyorderbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	supplyordercachemodel "tart-shop-manager/internal/entity/dtos/redis/supplyOrder"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type GetSupplyOrderStorage interface {
	GetSupplyOrder(ctx context.Context, cond map[string]interface{}) (*supplyordermodel.SupplyOrder, error)
}

type GetSupplyOrderCache interface {
	SaveSupplyOrder(ctx context.Context, data interface{}, morekeys ...string) error
	GetSupplyOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*supplyordermodel.SupplyOrder, error)
}

type getSupplyOrderBusiness struct {
	store GetSupplyOrderStorage
	cache GetSupplyOrderCache
}

func NewGetSupplyOrderBiz(store GetSupplyOrderStorage, cache GetSupplyOrderCache) *getSupplyOrderBusiness {
	return &getSupplyOrderBusiness{store, cache}
}

func (biz *getSupplyOrderBusiness) GetSupplyOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*supplyordermodel.SupplyOrder, error) {

	record, err := biz.cache.GetSupplyOrder(ctx, cond)
	if err != nil {
		return nil, common.ErrCannotGetEntity(supplyordermodel.EntityName, err)
	}

	if record != nil {
		return record, nil
	}

	record, err = biz.store.GetSupplyOrder(ctx, cond)

	if err != nil {
		return nil, common.ErrCannotGetEntity(supplyordermodel.EntityName, err)
	}

	if record != nil {

		var pagging paggingcommon.Paging
		pagging.Process()

		var createSupplyOrder = supplyordercachemodel.ToCreateSupplyOrderCache(*record)

		// Generate cache key
		key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
			EntityName: supplyordermodel.EntityName,
			Cond:       cond,
			Paging:     pagging,
			Filter:     commonfilter.Filter{},
			MoreKeys:   morekeys,
		})
		if err != nil {
			return nil, common.ErrCannotGenerateKey(supplyordermodel.EntityName, err)
		}

		if err := biz.cache.SaveSupplyOrder(ctx, createSupplyOrder, key); err != nil {
			return nil, common.ErrCannotCreateEntity(supplyordermodel.EntityName, err)
		}

	}

	return record, nil
}
