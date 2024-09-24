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

func (r *rdbStorage) GetRole(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*rolemodel.Role, error) {

	var paging paggingcommon.Paging
	paging.Process()

	// Generate cache key
	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: rolemodel.EntityName,
		Cond:       cond,
		Paging:     paging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return nil, common.ErrCannotGenerateKey(rolemodel.EntityName, err)
	}

	encryptedRecord, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	// Giải mã dữ liệu
	decryptedData, err := cacheutil.Decrypt([]byte(encryptedRecord))
	if err != nil {
		return nil, common.ErrDB(err)
	}

	var role rolemodel.Role

	if err := json.Unmarshal(decryptedData, &role); err != nil {
		return nil, common.ErrDB(err)
	}

	return &role, nil

}
