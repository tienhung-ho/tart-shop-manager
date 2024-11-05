package orderbusiness

import (
	"context"
	"errors"
	"tart-shop-manager/internal/common"
	ordermodel "tart-shop-manager/internal/entity/dtos/sql/order"
)

type DeleteOrderStorage interface {
	DeleteOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
	GetOrder(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*ordermodel.Order, error)
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type DeleteOrderItemStorage interface {
	DeleteOrderItem(ctx context.Context, cond map[string]interface{}, morekeys ...string) error
}

type DeleteOrderCache interface {
	DeleteListCache(ctx context.Context, entityName string) error
}

type deleteOrderBusiness struct {
	store DeleteOrderStorage
	cache DeleteOrderCache

	orderItemStorage DeleteOrderItemStorage
}

func NewDeleteOrderBiz(store DeleteOrderStorage, orderItemStorage DeleteOrderItemStorage) *deleteOrderBusiness {
	return &deleteOrderBusiness{store: store, orderItemStorage: orderItemStorage}
}

func (biz *deleteOrderBusiness) DeleteOrder(ctx context.Context,
	cond map[string]interface{}, morekeys ...string) error {

	record, err := biz.store.GetOrder(ctx, cond, morekeys...)

	if err != nil {
		if errors.Is(err, common.RecordNotFound) {
			return common.ErrNotFoundEntity(ordermodel.EntityName, err)
		}

		return common.ErrCannotGetEntity(ordermodel.EntityName, err)
	}

	err = biz.store.Transaction(ctx, func(txCtx context.Context) error {

		if err := biz.orderItemStorage.DeleteOrderItem(txCtx, map[string]interface{}{"order_id": record.OrderID}); err != nil {
			return common.ErrCannotDeleteEntity(ordermodel.EntityNameOrderItem, err)
		}

		if err := biz.store.DeleteOrder(txCtx, map[string]interface{}{"order_id": record.OrderID}, morekeys...); err != nil {
			return common.ErrCannotDeleteEntity(ordermodel.EntityName, err)
		}

		if err := biz.cache.DeleteListCache(txCtx, ordermodel.EntityName); err != nil {
			return common.ErrCannotDeleteEntity(ordermodel.EntityName, err)
		}

		return nil

	})

	return nil
}
