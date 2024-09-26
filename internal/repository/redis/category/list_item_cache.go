package categorycache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]categorymodel.Category, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var categories []categorymodel.Category
	if err := json.Unmarshal([]byte(record), &categories); err != nil {
		return nil, common.ErrDB(err)
	}
	return categories, nil
}
