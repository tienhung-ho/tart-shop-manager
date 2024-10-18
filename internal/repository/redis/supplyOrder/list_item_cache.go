package supplyordercache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	supplyordermodel "tart-shop-manager/internal/entity/dtos/sql/supplyOrder"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]supplyordermodel.SupplyOrder, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var supplyOrder []supplyordermodel.SupplyOrder
	if err := json.Unmarshal([]byte(record), &supplyOrder); err != nil {
		return nil, common.ErrDB(err)
	}
	return supplyOrder, nil
}
