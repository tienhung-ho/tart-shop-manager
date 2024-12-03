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

func (r *rdbStorage) GetProduct(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*productmodel.Product, error) {
	var paging paggingcommon.Paging
	paging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: productmodel.EntityName,
		Cond:       cond,
		Paging:     paging,
		Filter:     commonfilter.Filter{},
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

	var product productmodel.Product

	if err := json.Unmarshal([]byte(record), &product); err != nil {
		return nil, common.ErrDB(err)
	}

	return &product, nil
}

func (r *rdbStorage) GetPaging(ctx context.Context, key string) (*paggingcommon.Paging, error) {
	// Lấy giá trị từ Redis
	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("paging not found")
		}
		return nil, common.ErrDB(err)
	}

	// Unmarshal dữ liệu JSON thành struct `Paging`
	var paging paggingcommon.Paging
	if err := json.Unmarshal([]byte(result), &paging); err != nil {
		return nil, common.ErrDB(err)
	}

	return &paging, nil
}

func (r *rdbStorage) GetFilter(ctx context.Context, key string) (*commonfilter.Filter, error) {
	// Lấy giá trị từ Redis
	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("filter not found")
		}
		return nil, common.ErrDB(err)
	}

	// Unmarshal dữ liệu JSON thành struct `Filter`
	var filter commonfilter.Filter
	if err := json.Unmarshal([]byte(result), &filter); err != nil {
		return nil, common.ErrDB(err)
	}

	return &filter, nil
}
