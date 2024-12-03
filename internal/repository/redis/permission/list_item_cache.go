package permissioncachestorage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
	cacheutil "tart-shop-manager/internal/util/cache"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]permissionmodel.Permission, error) {
	encryptedRecord, err := r.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil // cache misss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	// Giải mã dữ liệu
	decryptedData, err := cacheutil.Decrypt([]byte(encryptedRecord))
	var permissions []permissionmodel.Permission
	if err := json.Unmarshal([]byte(decryptedData), &permissions); err != nil {
		return nil, common.ErrDB(err)
	}
	return permissions, nil
}
