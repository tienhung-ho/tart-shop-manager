package ingredientcache

import (
	"context"
	"fmt"
	"log"
	"tart-shop-manager/internal/common"
)

func (r *rdbStorage) DeleteListCache(ctx context.Context, entityName string) error {
	var cursor uint64
	var keys []string
	var err error

	// Sử dụng SCAN để tìm các keys bắt đầu với "cache:list:"
	for {
		var batch []string
		batch, cursor, err = r.rdb.Scan(ctx, cursor, fmt.Sprintf("cache:List:%s:*", entityName), 100).Result()
		if err != nil {
			return common.ErrDB(err)
		}
		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}

	log.Print(keys)
	if len(keys) > 0 {
		if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
			return common.ErrDB(err)
		}
	}

	return nil
}
