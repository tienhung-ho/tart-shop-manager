package accountrdbstorage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]accountmodel.Account, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var records []accountmodel.Account
	if err := json.Unmarshal([]byte(record), &records); err != nil {
		return nil, common.ErrDB(err)
	}
	return records, nil
}
