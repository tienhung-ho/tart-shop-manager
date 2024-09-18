package ingredientcache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	ingredientmodel "tart-shop-manager/internal/entity/dtos/sql/ingredient"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) GetIngredient(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ingredientmodel.Ingredient, error) {
	var paging paggingcommon.Paging
	paging.Process()

	key := cacheutil.GenerateKey(categorymodel.EntityName, cond, paging, commonfilter.Filter{})

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var ingredient ingredientmodel.Ingredient

	if err := json.Unmarshal([]byte(record), &ingredient); err != nil {
		return nil, common.ErrDB(err)
	}

	return &ingredient, nil
}
