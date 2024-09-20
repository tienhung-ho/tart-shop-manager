package rolecache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) ListItemRole(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]rolemodel.Role, error) {
	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: ordermodel.EntityName,
		Cond:       cond,
		Paging:     *paging,
		Filter:     *filter,
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(ordermodel.EntityName, err)
	}

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var roles []rolemodel.Role

	if err := json.Unmarshal([]byte(record), &roles); err != nil {
		return nil, common.ErrDB(err)
	}

	return roles, nil
}
