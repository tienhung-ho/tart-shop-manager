package paggingcommon

import (
	"fmt"
	"strings"
)

type Paging struct {
	Page  int    `json:"page" form:"page"`
	Limit int    `json:"limit" form:"limit"`
	Total int64  `json:"total" form:"-"`
	Sort  string `json:"sort" form:"sort"`
}

func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 9999
	}
}

func (p *Paging) ParseSortFields(sortParam string, allowedSortFields map[string]bool) ([]string, error) {
	if sortParam == "" {
		return nil, nil
	}

	sortFields := strings.Split(sortParam, ",")
	var validSorts []string

	for _, field := range sortFields {
		parts := strings.Fields(strings.TrimSpace(field))
		if len(parts) == 0 {
			continue
		}

		fieldName := parts[0]
		direction := "asc"

		if len(parts) > 1 {
			dir := strings.ToLower(parts[1])
			if dir == "asc" || dir == "desc" {
				direction = dir
			} else {
				return nil, fmt.Errorf("invalid sort direction: %s", dir)
			}
		}

		// Validate field name
		if !allowedSortFields[fieldName] {
			return nil, fmt.Errorf("invalid sort field: %s", fieldName)
		}

		validSorts = append(validSorts, fmt.Sprintf("%s %s", fieldName, direction))
	}

	return validSorts, nil
}
