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

func (r *rdbStorage) GetCategory(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*categorymodel.Category, error) {
	var paging paggingcommon.Paging
	paging.Process()

	key := cacheutil.GenerateKey(categorymodel.EntityName, cond, paging, commonfilter.Filter{})

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var category categorymodel.Category

	if err := json.Unmarshal([]byte(record), &category); err != nil {
		return nil, common.ErrDB(err)
	}

	return &category, nil
}
