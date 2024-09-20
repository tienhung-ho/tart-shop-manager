package accountrdbstorage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	"tart-shop-manager/internal/entity/dtos/sql/account"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error) {

	var paging paggingcommon.Paging
	paging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: accountmodel.EntityName,
		Cond:       cond,
		Paging:     paging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(accountmodel.EntityName, err)
	}

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var account accountmodel.Account

	if err := json.Unmarshal([]byte(record), &account); err != nil {
		return nil, common.ErrDB(err)
	}

	return &account, nil
}
