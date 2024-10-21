package stockbatchcache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) GetStockBatch(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*stockbatchmodel.StockBatch, error) {
	var paging paggingcommon.Paging
	paging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: stockbatchmodel.EntityName,
		Cond:       cond,
		Paging:     paging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(stockbatchmodel.EntityName, err)
	}

	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var stockBatch stockbatchmodel.StockBatch

	if err := json.Unmarshal([]byte(record), &stockBatch); err != nil {
		return nil, common.ErrDB(err)
	}

	return &stockBatch, nil
}
