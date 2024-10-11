package recipecache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]recipemodel.Recipe, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var recipes []recipemodel.Recipe
	if err := json.Unmarshal([]byte(record), &recipes); err != nil {
		return nil, common.ErrDB(err)
	}
	return recipes, nil
}
