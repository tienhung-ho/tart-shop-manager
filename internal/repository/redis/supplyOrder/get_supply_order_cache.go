package supplyordercache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) GetSupplyOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*supplyordermodel.SupplyOrder, error) {
	var paging paggingcommon.Paging
	paging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: supplyordermodel.EntityName,
		Cond:       cond,
		Paging:     paging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(supplyordermodel.EntityName, err)
	}

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var supplyOrder supplyordermodel.SupplyOrder

	if err := json.Unmarshal([]byte(record), &supplyOrder); err != nil {
		return nil, common.ErrDB(err)
	}

	return &supplyOrder, nil
}
