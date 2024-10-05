package recipecache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	recipemodel "tart-shop-manager/internal/entity/dtos/sql/recipe"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) GetRecipe(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*recipemodel.Recipe, error) {
	var paging paggingcommon.Paging
	paging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: recipemodel.EntityName,
		Cond:       cond,
		Paging:     paging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(recipemodel.EntityName, err)
	}

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var recipe recipemodel.Recipe

	if err := json.Unmarshal([]byte(record), &recipe); err != nil {
		return nil, common.ErrDB(err)
	}

	return &recipe, nil
}
