package databaseutil

func PrefixConditionKeys(cond map[string]interface{}, alias string) map[string]interface{} {
	prefixedCond := make(map[string]interface{})
	for key, value := range cond {
		prefixedCond[alias+"."+key] = value
	}
	return prefixedCond
}
