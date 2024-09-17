package ordercache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error) {

	var paging paggingcommon.Paging
	paging.Process()

	key := cacheutil.GenerateKey(ordermodel.EntityName, cond, paging, commonfilter.Filter{})

	record, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var order ordermodel.Order

	if err := json.Unmarshal([]byte(record), &order); err != nil {
		return nil, common.ErrDB(err)
	}

	return &order, nil

}
