package databaseutil

func PrefixConditionKeys(cond map[string]interface{}, alias string) map[string]interface{} {
	prefixedCond := make(map[string]interface{})
	for key, value := range cond {
		prefixedCond[alias+"."+key] = value
	}
	return prefixedCond
}

func Difference(a, b []uint64) []uint64 {
	m := make(map[uint64]bool)
	for _, item := range b {
		m[item] = true
	}

	var diff []uint64
	for _, item := range a {
		if !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}
