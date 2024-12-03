package permissioncachestorage

import (
	"context"
	"encoding/json"
	"errors"
	"tart-shop-manager/internal/common"
	permissionmodel "tart-shop-manager/internal/entity/dtos/sql/permission"
	cacheutil "tart-shop-manager/internal/util/cache"
	"time"
)

func (r *rdbStorage) SavePermission(ctx context.Context, data interface{}, morekeys ...string) error {
	if len(morekeys) == 0 {
		return common.ErrDB(errors.New("missing cache key"))
	}

	key := morekeys[0]
	var record []byte
	var err error

	switch v := data.(type) {
	case []permissionmodel.Permission:
		record, err = json.Marshal(v)
		if err != nil {
			return common.ErrDB(err)
		}
	default:
		return errors.New("unsupported data type")
	}

	// Mã hóa dữ liệu
	encryptedData, err := cacheutil.Encrypt(record)
	if err != nil {
		return common.ErrDB(err)
	}

	if err := r.rdb.Set(ctx, key, encryptedData, 10*time.Minute).Err(); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
