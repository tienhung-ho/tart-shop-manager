package categorycache

import (
	"context"
	"encoding/json"
	"errors"
	"tart-shop-manager/internal/common"
	categorycachemodel "tart-shop-manager/internal/entity/dtos/redis/category"
	categorymodel "tart-shop-manager/internal/entity/dtos/sql/category"
	"time"
)

func (r *rdbStorage) SaveCategory(ctx context.Context, data interface{}, morekeys ...string) error {
	if len(morekeys) == 0 {
		return common.ErrDB(errors.New("missing cache key"))
	}

	key := morekeys[0]
	var record []byte
	var err error

	switch v := data.(type) {
	case *categorycachemodel.CreateCategory:

		record, err = json.Marshal(v)

		if err != nil {
			return common.ErrDB(err)
		}
	case []categorymodel.Category:
		record, err = json.Marshal(v)
		if err != nil {
			return common.ErrDB(err)
		}
	default:
		return errors.New("unsupported data type")
	}
	if err := r.rdb.Set(ctx, key, string(record), 20*time.Minute).Err(); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
