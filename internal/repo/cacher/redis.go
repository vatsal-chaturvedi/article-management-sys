package cacher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/vatsal-chaturvedi/article-management-sys/internal/config"
	"time"
)

//go:generate mockgen --build_flags=--mod=mod --destination=./../../../pkg/mock/mock_cacher.go --package=mock github.com/vatsal-chaturvedi/article-management-sys/internal/repo/cacher CacherI
type CacherI interface {
	Get(string) ([]byte, error)
	Set(string, interface{}, time.Duration) error
}
type cache struct {
	rdb *redis.Client
}

func NewCacher(cacheSvc config.CacheSvc) CacherI {
	return &cache{
		rdb: cacheSvc.Rdb,
	}
}

func (c cache) Get(key string) ([]byte, error) {
	data, err := c.rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	return data, err
}
func (c cache) Set(key string, value interface{}, expiry time.Duration) error {
	err := c.rdb.Set(context.Background(), key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}
