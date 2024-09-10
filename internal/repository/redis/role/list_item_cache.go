package rolecache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) ListItemRole(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging, filter *commonfilter.Filter, morekeys ...string) ([]rolemodel.Role, error) {
	key := cacheutil.GenerateKey(rolemodel.EntityName, cond, *paging, *filter)

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
