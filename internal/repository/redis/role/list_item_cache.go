package rolecache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]rolemodel.Role, error) {
	
	encryptedRecord, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	// Giải mã dữ liệu
	decryptedData, err := cacheutil.Decrypt([]byte(encryptedRecord))

	var roles []rolemodel.Role
	if err := json.Unmarshal([]byte(decryptedData), &roles); err != nil {
		return nil, common.ErrDB(err)
	}
	return roles, nil
}
