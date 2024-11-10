package commonfilter

import (
	"tart-shop-manager/internal/common"
)

type Filter struct {
	Status     string             `json:"status,omitempty" form:"status"`
	Search     string             `json:"search,omitempty" form:"search"`
	MinPrice   float64            `json:"min_price,omitempty" form:"min_price"`
	MaxPrice   float64            `json:"max_price,omitempty" form:"max_price"`
	CategoryID uint64             `json:"category_id,omitempty" form:"category_id"`
	IDs        []uint64           `json:"ids,omitempty"`
	StartDate  *common.CustomDate `json:"start_date,omitempty" form:"start_date"`
	EndDate    *common.CustomDate `json:"end_date,omitempty" form:"end_date"`
	Ingredient *uint64            `json:"ingredient,omitempty" form:"ingredient"`
	OrderDate
	Recipe
	Product
	ReportOrderDate
}

type OrderDate struct {
	ExpirationDate      *common.CustomDate `json:"expiration_date,omitempty" form:"expiration_date"`
	ReceivedDate        *common.CustomDate `json:"received_date,omitempty" form:"received_date"`
	StartExpirationDate *common.CustomDate `json:"start_expiration_date,omitempty" form:"start_expiration_date"`
	EndExpirationDate   *common.CustomDate `json:"end_expiration_date,omitempty" form:"end_expiration_date"`
	StartReceivedDate   *common.CustomDate `json:"start_received_date,omitempty" form:"start_received_date"`
	EndReceivedDate     *common.CustomDate `json:"end_received_date,omitempty" form:"end_received_date"`
}

type ReportOrderDate struct {
	InDate *common.CustomDate `json:"in_date,omitempty" form:"in_date"`
}

type Recipe struct {
	ProductIDs []uint64 `json:"product_ids,omitempty"`
	Sizes      []string `json:"sizes,omitempty"`
}

type Product struct {
	//ProductIDs []uint64 `json:"product_ids,omitempty"`
	Name string `json:"product_name,omitempty"`
}
