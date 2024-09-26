package rolecache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"tart-shop-manager/internal/common"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
)

func (r *rdbStorage) ListItem(ctx context.Context, key string) ([]rolemodel.Role, error) {
	record, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, common.ErrDB(err)
	}

	var roles []rolemodel.Role
	if err := json.Unmarshal([]byte(record), &roles); err != nil {
		return nil, common.ErrDB(err)
	}
	return roles, nil
}
