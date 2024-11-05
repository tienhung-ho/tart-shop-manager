package ordercache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]ordermodel.Order, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var orders []ordermodel.Order
	if err := json.Unmarshal([]byte(record), &orders); err != nil {
		return nil, common.ErrDB(err)
	}
	return orders, nil
}
