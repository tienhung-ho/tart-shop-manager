package reportbusiness

import (
	"context"
	"tart-shop-manager/internal/common"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
	reportmodel "tart-shop-manager/internal/entity/dtos/sql/report"
)

type OrderReportStorage interface {
	ReportOrder(ctx context.Context, cond map[string]interface{}, pagging *paggingcommon.Paging,
		filter *commonfilter.Filter, morekeys ...string) ([]reportmodel.Order, error)
}

type orderReportBusiness struct {
	store OrderReportStorage
}

func NewOrderReportBiz(store OrderReportStorage) *orderReportBusiness {
	return &orderReportBusiness{store}
}

func (biz *orderReportBusiness) ReportOrder(ctx context.Context, cond map[string]interface{},
	pagging *paggingcommon.Paging,
	filter *commonfilter.Filter, morekeys ...string) (*reportmodel.RevenueReport, error) {

	records, err := biz.store.ReportOrder(ctx, cond, pagging, filter, morekeys...)
	if err != nil {
		return nil, common.ErrCannotGetReport(reportmodel.OrderReportName, err)
	}
	// Chuyển đổi records thành RevenueReport
	report, err := biz.generateRevenueReport(records, filter)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (biz *orderReportBusiness) generateRevenueReport(records []reportmodel.Order,
	filter *commonfilter.Filter) (*reportmodel.RevenueReport, error) {

	var totalRevenue float64 = 0
	var totalCost float64 = 0
	var orderSummaries []reportmodel.OrderSummary

	for _, record := range records {
		totalRevenue += record.TotalAmount

		var orderItemSummary []reportmodel.OrderItemSummary
		for _, item := range record.OrderItems {
			totalCost += item.Recipe.Cost * float64(item.Quantity)
			orderItemSummary = append(orderItemSummary, reportmodel.OrderItemSummary{
				ProductID: item.Recipe.ProductID,
				Name:      item.Recipe.Product.Name,
				Price:     item.Recipe.Product.Price,
				ImageURL:  item.Recipe.Product.ImageURL,
				RecipeID:  item.RecipeID,
				Size:      item.Recipe.Size,
				Cost:      item.Recipe.Cost,
			})
		}
		orderSummaries = append(orderSummaries, reportmodel.OrderSummary{
			OrderID:     record.OrderID,
			OrderDate:   &common.CustomDate{Time: record.CreatedAt},
			TotalAmount: record.TotalAmount,
			Items:       orderItemSummary,
		})

	}

	report := &reportmodel.RevenueReport{
		StartDate:    filter.StartDate, // Bạn cần truyền StartDate và EndDate vào filter hoặc cond
		EndDate:      filter.EndDate,
		TotalRevenue: totalRevenue,
		TotalCost:    totalCost,
		TotalOrders:  len(records),
		Orders:       orderSummaries,
	}

	return report, nil
}
