package productcache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) ListItem(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]productmodel.Product, error) {

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: productmodel.EntityName,
		Cond:       cond,
		Paging:     *pagging,
		Filter:     *filter,
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(productmodel.EntityName, err)
	}

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var products []productmodel.Product

	if err := json.Unmarshal([]byte(record), &products); err != nil {
		return nil, common.ErrDB(err)
	}
	return nil, nil
}
