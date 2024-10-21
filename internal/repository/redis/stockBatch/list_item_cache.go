package stockbatchcache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	stockbatchmodel "tart-shop-manager/internal/entity/dtos/sql/stockbatch"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]stockbatchmodel.StockBatch, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var stockBatches []stockbatchmodel.StockBatch
	if err := json.Unmarshal([]byte(record), &stockBatches); err != nil {
		return nil, common.ErrDB(err)
	}
	return stockBatches, nil
}
