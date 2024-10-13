package productcache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"tart-shop-manager/internal/common"
	productcachemodel "tart-shop-manager/internal/entity/dtos/redis/product"
	productmodel "tart-shop-manager/internal/entity/dtos/sql/product"
	"time"
)

func (r *rdbStorage) SaveProduct(ctx context.Context, data interface{}, morekeys ...string) error {
	if len(morekeys) == 0 {
		return common.ErrDB(errors.New("missing cache key"))
	}

	key := morekeys[0]
	var record []byte
	var err error

	log.Print(key)

	switch v := data.(type) {
	case *productcachemodel.CreateProduct:

		record, err = json.Marshal(v)

		if err != nil {
			return common.ErrDB(err)
		}
	case []productmodel.Product:
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
