package ingredientcache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]ingredientmodel.Ingredient, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var ingredients []ingredientmodel.Ingredient
	if err := json.Unmarshal([]byte(record), &ingredients); err != nil {
		return nil, common.ErrDB(err)
	}
	return ingredients, nil
}
