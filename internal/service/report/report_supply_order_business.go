package reportbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	reportmodel "tart-shop-manager/internal/entity/dtos/sql/report"
)

type ReportSupplyOrderStorage interface {
	ReportSupplyOrder(ctx context.Context, cond map[string]interface{}, paging *paggingcommon.Paging,
		filter *commonfilter.Filter, morekeys ...string) ([]reportmodel.SupplyOrder, error)
}

type reportSupplyOrderBusiness struct {
	store ReportSupplyOrderStorage
}

func NewReportSupplyOrderBusiness(store ReportSupplyOrderStorage) *reportSupplyOrderBusiness {
	return &reportSupplyOrderBusiness{store: store}
}

func (biz *reportSupplyOrderBusiness) ReportSupplyOrder(ctx context.Context,
	cond map[string]interface{}, paging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) (*reportmodel.SupplyReport, error) {

	records, err := biz.store.ReportSupplyOrder(ctx, cond, paging, filter, morekeys...)
	if err != nil {
		return nil, common.ErrCannotGetReport(reportmodel.SupplyOrderReportName, err)
	}

	report, err := biz.generateSupplyOrderReport(records, filter)

	if err != nil {
		return nil, common.ErrCannotGetReport(reportmodel.SupplyOrderReportName, err)
	}

	return report, nil
}

func (biz *reportSupplyOrderBusiness) generateSupplyOrderReport(records []reportmodel.SupplyOrder,
	filter *commonfilter.Filter) (*reportmodel.SupplyReport, error) {

	var totalSupplyCost float64 = 0
	var supplyOrderSummaries []reportmodel.SupplyOrderSummary

	for _, record := range records {
		totalSupplyCost += record.TotalAmount

		var supplyOrderItemSummary []reportmodel.SupplyOrderItemSummary
		for _, item := range record.SupplyOrderItems {
			supplyOrderItemSummary = append(supplyOrderItemSummary, reportmodel.SupplyOrderItemSummary{
				IngredientID:   item.IngredientID,
				IngredientName: item.Ingredient.Name,
				Price:          item.Price,
				Quantity:       item.Quantity,
				Unit:           item.Unit,
				TotalCost:      item.Quantity * item.Price,
			})
		}

		supplyOrderSummaries = append(supplyOrderSummaries, reportmodel.SupplyOrderSummary{
			SupplyOrderID: record.SupplyOrderID,
			OrderDate:     &common.CustomDate{Time: record.OrderDate},
			SupplierID:    record.SupplyOrderID,
			SupplierName:  record.Supplier.Name,
			TotalAmount:   record.TotalAmount,
			Items:         supplyOrderItemSummary,
			ContactInfo:   record.Supplier.ContactInfo,
		})
	}

	report := &reportmodel.SupplyReport{
		StartDate:         filter.StartDate,
		EndDate:           filter.EndDate,
		TotalSupplyCost:   totalSupplyCost,
		TotalSupplyOrders: len(records),
		SupplyOrders:      supplyOrderSummaries,
	}

	return report, nil
}
