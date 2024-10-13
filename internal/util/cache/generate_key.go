package cacheutil

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sort"
	commonfilter "tart-shop-manager/internal/common/filter"
	paggingcommon "tart-shop-manager/internal/common/paging"
)

type CacheParams struct {
	EntityName string
	Cond       map[string]interface{}
	Paging     paggingcommon.Paging
	Filter     commonfilter.Filter
	MoreKeys   []string
	KeyType    string // Thêm trường này
}

func GenerateKey(params CacheParams) (string, error) {
	// Serialize cond with sorted keys
	condBytes, err := marshalMap(params.Cond)
	if err != nil {
		return "", fmt.Errorf("error marshaling cond: %w", err)
	}

	// Create a canonical representation
	keyData := struct {
		EntityName string               `json:"entity_name"`
		Cond       json.RawMessage      `json:"cond"`
		Paging     paggingcommon.Paging `json:"paging"`
		Filter     commonfilter.Filter  `json:"filter"`
		MoreKeys   []string             `json:"more_keys"`
		KeyType    string               `json:"key_type"`
	}{
		EntityName: params.EntityName,
		Cond:       condBytes,
		Paging:     params.Paging,
		Filter:     params.Filter,
		MoreKeys:   params.MoreKeys,
		KeyType:    params.KeyType, // Bao gồm KeyType
	}

	jsonBytes, err := json.Marshal(keyData)
	if err != nil {
		return "", fmt.Errorf("error marshaling key data: %w", err)
	}

	hash := md5.Sum(jsonBytes)
	return fmt.Sprintf("cache:%s:%x", params.KeyType, hash), nil // Thêm KeyType vào prefix
}

func marshalMap(m map[string]interface{}) ([]byte, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	orderedMap := make([]struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}, 0, len(m))

	for _, k := range keys {
		orderedMap = append(orderedMap, struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
		}{
			Key:   k,
			Value: m[k],
		})
	}

	return json.Marshal(orderedMap)
}
