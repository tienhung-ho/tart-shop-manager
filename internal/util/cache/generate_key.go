package cacheutil

import (
	"fmt"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
)

func GenerateKey(entityName string, cond map[string]interface{}, paging paggingcommon.Paging,
	filter commonfilter.Filter, morekeys ...string) string {

	key := fmt.Sprintf("%s:", entityName)
	for k, v := range cond {
		key += fmt.Sprintf("%s=%v:", k, v)
	}

	key += fmt.Sprintf("%s=%v:", "page", paging.Page)
	key += fmt.Sprintf("%s=%v:", "limit", paging.Limit)

	key += fmt.Sprintf("%s=%v:", "filter", filter.Status)

	return key
}
