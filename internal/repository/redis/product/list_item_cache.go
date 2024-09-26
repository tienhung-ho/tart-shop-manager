package productcache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]productmodel.Product, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var products []productmodel.Product
	if err := json.Unmarshal([]byte(record), &products); err != nil {
		return nil, common.ErrDB(err)
	}
	return products, nil
}
