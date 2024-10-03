package accountbusiness

import (
	"context"
	"sync"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	accountmodel "tart-shop-manager/internal/entity/dtos/sql/account"
	imagemodel "tart-shop-manager/internal/entity/dtos/sql/image"
	rolemodel "tart-shop-manager/internal/entity/dtos/sql/role"
	casbinbusiness "tart-shop-manager/internal/service/policies"
	rolebusiness "tart-shop-manager/internal/service/role"
	cacheutil "tart-shop-manager/internal/util/cache"
)

type UpdateAccountStorage interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	UpdateAccount(ctx context.Context, cond map[string]interface{}, account *accountmodel.UpdateAccount, morekeys ...string) (*accountmodel.Account, error)
}

type UpdateAccountCache interface {
	GetAccount(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*accountmodel.Account, error)
	DeleteAccount(ctx context.Context, morekeys ...string) error
}

type UpdateImage interface {
	DeleteImage(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	UpdateImage(ctx context.Context, cond map[string]interface{}, data *imagemodel.UpdateImage) error
	ListItem(ctx context.Context, cond map[string]interface{}) ([]imagemodel.Image, error)
}

type updateAccountBusiness struct {
	store     UpdateAccountStorage
	roleStore rolebusiness.GetRoleStorage
	cache     UpdateAccountCache
	image     UpdateImage
	auth      casbinbusiness.Authorization
}

func NewUpdateAccount(store UpdateAccountStorage, roleStore rolebusiness.GetRoleStorage, cache UpdateAccountCache, image UpdateImage, authorization casbinbusiness.Authorization) *updateAccountBusiness {
	return &updateAccountBusiness{store, roleStore, cache, image, authorization}
}

func (biz *updateAccountBusiness) UpdateAccount(ctx context.Context, cond map[string]interface{}, account *accountmodel.UpdateAccount, morekeys ...string) (*accountmodel.Account, error) {

	record, err := biz.store.GetAccount(ctx, cond)

	if err != nil {
		return nil, common.ErrNotFoundEntity(accountmodel.EntityName, err)
	}

	updatedRecord, err := biz.store.UpdateAccount(ctx, map[string]interface{}{"account_id": record.AccountID}, account, morekeys...)

	if err != nil {
		return nil, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	imageCond := map[string]interface{}{
		"account_id": record.AccountID,
	}

	oldImages, err := biz.image.ListItem(ctx, imageCond)

	if err != nil {
		return nil, common.ErrCannotListEntity(imagemodel.EntityName, err)
	}

	oldImageIDs := make(map[uint64]bool)
	for _, img := range oldImages {
		oldImageIDs[img.ImageID] = true
	}

	// 4. Xác định hình ảnh cần thêm và cần xóa
	newImageIDs := make(map[uint64]bool)
	for _, imgID := range account.Images {
		newImageIDs[imgID.ImageID] = true
	}

	var imagesToAdd []uint64
	var imagesToRemove []uint64

	for imgID := range newImageIDs {
		if !oldImageIDs[imgID] {
			imagesToAdd = append(imagesToAdd, imgID)
		}
	}

	for imgID := range oldImageIDs {
		if !newImageIDs[imgID] {
			imagesToRemove = append(imagesToRemove, imgID)
		}
	}

	// 5. Cập nhật `ProductID` cho hình ảnh cần thêm và xóa sử dụng goroutines
	var wg sync.WaitGroup
	var updateErr error
	var mu sync.Mutex

	// Giới hạn số lượng goroutines chạy đồng thời để tránh làm quá tải hệ thống
	sem := make(chan struct{}, 10) // Giới hạn 10 goroutines

	// Thêm hình ảnh mới (Cập nhật `ProductID` cho các hình ảnh này)
	for _, imgID := range imagesToAdd {
		wg.Add(1)
		sem <- struct{}{} // Acquire a token
		go func(id uint64) {
			defer wg.Done()
			defer func() { <-sem }() // Release the token

			err := biz.image.UpdateImage(ctx, map[string]interface{}{"image_id": id}, &imagemodel.UpdateImage{
				AccountID: &record.AccountID,
			})
			if err != nil {
				mu.Lock()
				updateErr = err
				mu.Unlock()
			}
		}(imgID)
	}

	// Xóa liên kết hình ảnh cũ (Đặt `ProductID` về NULL)
	for _, imgID := range imagesToRemove {
		wg.Add(1)
		sem <- struct{}{}
		go func(id uint64) {
			defer wg.Done()
			defer func() { <-sem }()

			err := biz.image.DeleteImage(ctx, map[string]interface{}{"image_id": id})
			if err != nil {
				mu.Lock()
				updateErr = err
				mu.Unlock()
			}
		}(imgID)
	}

	wg.Wait()

	if updateErr != nil {
		// Xử lý lỗi (rollback transaction nếu cần)
		return nil, common.ErrCannotUpdateEntity("Image", updateErr)
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
		return nil, common.ErrCannotGenerateKey(accountmodel.EntityName, err)
	}

	if err := biz.cache.DeleteAccount(ctx, key); err != nil {
		return nil, common.ErrCannotUpdateEntity(accountmodel.EntityName, err)
	}

	if record.RoleID != updatedRecord.RoleID {
		// Xóa vai trò cũ của người dùng
		if err := biz.auth.RemoveUserFromAllRoles(ctx, record.Email); err != nil {
			return nil, common.ErrCannotDeleteEntity("user roles", err)
		}

		role, err := biz.roleStore.GetRole(ctx, map[string]interface{}{"role_id": updatedRecord.RoleID})

		if err != nil {
			return nil, common.ErrNotFoundEntity(rolemodel.EntityName, err)
		}

		if err := biz.auth.AddUserToRole(ctx, record.Email, role.Name); err != nil {
			return nil, common.ErrCannotCreateEntity("user roles", err)
		}

	}

	return updatedRecord, nil
}
