package supplyordercache

import (
	"context"
	"tart-shop-manager/internal/common"
)

func (r *rdbStorage) DeleteSupplyOrder(ctx context.Context, morekeys ...string) error {
	key := morekeys[0]
	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
