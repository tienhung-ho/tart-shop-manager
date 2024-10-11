package commonfilter

type Filter struct {
	Status     string   `json:"status" form:"status"`
	Search     string   `json:"search" form:"search"`
	MinPrice   float64  `json:"min_price" form:"min_price"`
	MaxPrice   float64  `json:"max_price" form:"max_price"`
	CategoryID uint64   `json:"category_id" form:"category_id"`
	IDs        []uint64 `json:"ids,omitempty"`
}
