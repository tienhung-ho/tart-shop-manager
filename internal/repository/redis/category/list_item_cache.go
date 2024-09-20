package categorycache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) ListItem(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]categorymodel.Category, error) {

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: categorymodel.EntityName,
		Cond:       cond,
		Paging:     *paging,
		Filter:     *filter,
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(categorymodel.EntityName, err)
	}

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var categories []categorymodel.Category

	if err := json.Unmarshal([]byte(record), &categories); err != nil {
		return nil, common.ErrDB(err)
	}

	return categories, nil
}
