package suppliercache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	suppliermodel "tart-shop-manager/internal/entity/dtos/sql/supplier"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]suppliermodel.Supplier, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var suppliers []suppliermodel.Supplier
	if err := json.Unmarshal([]byte(record), &suppliers); err != nil {
		return nil, common.ErrDB(err)
	}
	return suppliers, nil
}
