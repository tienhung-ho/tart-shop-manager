package permissionstorage

import "context"

func (m *mysqlPermission) GetPermission(ctx context.Context, cond map[string]interface{}, morekeys ...string)
