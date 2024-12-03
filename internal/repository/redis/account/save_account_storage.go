package accountrdbstorage

import (
	"context"
	"encoding/json"
	"errors"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountrdbmodel "tart-shop-manager/internal/entity/dtos/redis/account"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	"time"
)

func (r *rdbStorage) SaveAccount(ctx context.Context, data interface{}, morekeys ...string) error {
	if len(morekeys) == 0 {
		return common.ErrDB(errors.New("missing cache key"))
	}

	key := morekeys[0]
	var record []byte
	var err error

	switch v := data.(type) {
	case *accountrdbmodel.CreateAccountRdb:

		record, err = json.Marshal(v)

		if err != nil {
			return common.ErrDB(err)
		}
	case []accountmodel.Account:
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

func (r *rdbStorage) SavePaging(ctx context.Context, paging *paggingcommon.Paging, morekeys ...string) error {
	if len(morekeys) == 0 {
		return common.ErrDB(errors.New("missing cache key"))
	}

	key := morekeys[0]

	// Mã hóa `paging` thành JSON
	record, err := json.Marshal(paging)
	if err != nil {
		return common.ErrDB(err)
	}

	// Lưu vào Redis
	if err := r.rdb.Set(ctx, key, string(record), 20*time.Minute).Err(); err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (r *rdbStorage) SaveFilter(ctx context.Context, filter *commonfilter.Filter, morekeys ...string) error {
	if len(morekeys) == 0 {
		return common.ErrDB(errors.New("missing cache key"))
	}

	key := morekeys[0]

	// Mã hóa `filter` thành JSON
	record, err := json.Marshal(filter)
	if err != nil {
		return common.ErrDB(err)
	}

	// Lưu vào Redis
	if err := r.rdb.Set(ctx, key, string(record), 20*time.Minute).Err(); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
