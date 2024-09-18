package ingredientcache

import "github.com/redis/go-redis/v9"

type rdbStorage struct {
	rdb *redis.Client
}

func NewRdbStorage(rdb *redis.Client) *rdbStorage {
	return &rdbStorage{rdb: rdb}
}
