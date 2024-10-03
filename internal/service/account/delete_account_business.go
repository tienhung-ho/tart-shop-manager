package accountbusiness

import (
	"context"
	"sync"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type DeleteAccountStorage interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	DeleteAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type DeleteAccountCache interface {
	DeleteAccount(ctx context.Context, morekeys ...string) error
}

type DeleteImage interface {
	DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error
	ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error)
}

type deleteAccountBusiness struct {
	store DeleteAccountStorage
	cache DeleteAccountCache
	image DeleteImage
}

func NewDeleteAccountBiz(store DeleteAccountStorage, cache DeleteAccountCache, image DeleteImage) *deleteAccountBusiness {
	return &deleteAccountBusiness{store, cache, image}
}

func (biz *deleteAccountBusiness) DeleteAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetAccount(ctx, cond)

	if err != nil {
		return common.ErrNotFoundEntity(accountmodel.EntityName, err)
	}

	if err := biz.store.DeleteAccount(ctx, map[string]interface{}{"account_id": record.AccountID}, morekeys...); err != nil {
		return common.ErrCannotDeleteEntity(accountmodel.EntityName, err)
	}

	// 3. Lấy danh sách hình ảnh cũ liên kết với sản phẩm
	imageCond := map[string]interface{}{
		"account_id": record.AccountID,
	}

	oldImages, err := biz.image.ListItem(ctx, imageCond)

	if err != nil {
		return common.ErrCannotListEntity(imagemodel.EntityName, err)
	}

	var imagesToRemove []uint64
	for _, imgID := range oldImages {
		imagesToRemove = append(imagesToRemove, imgID.ImageID)
	}

	var wg sync.WaitGroup
	var deleteErr error
	var mu sync.Mutex

	// Giới hạn số lượng goroutines chạy đồng thời để tránh làm quá tải hệ thống
	sem := make(chan struct{}, 10) // Giới hạn 10 goroutines

	for _, imgID := range imagesToRemove {
		wg.Add(1)
		sem <- struct{}{}
		go func(id uint64) {
			defer wg.Done()
			defer func() { <-sem }()
			err := biz.image.DeleteImage(ctx, map[string]interface{}{"image_id": id})
			if err != nil {
				mu.Lock()
				deleteErr = err
				mu.Unlock()
			}
		}(imgID)
	}

	wg.Wait()

	if deleteErr != nil {
		// Xử lý lỗi (rollback transaction nếu cần)
		return common.ErrCannotUpdateEntity("Image", deleteErr)
	}

	var pagging paggingcommon.Paging
	pagging.Process()

	key, err := cacheutil.GenerateKey(cacheutil.CacheParams{
		EntityName: accountmodel.EntityName,
		Cond:       cond,
		Paging:     pagging,
		Filter:     commonfilter.Filter{},
		MoreKeys:   morekeys,
	})
	if err != nil {
		return common.ErrCannotGenerateKey(accountmodel.EntityName, err)
	}

	if err := biz.cache.DeleteAccount(ctx, key); err != nil {
		return common.ErrCannotDeleteEntity(accountmodel.EntityName, err)
	}

	return nil
}
