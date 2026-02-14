package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{client: client}
}

// Get 获取缓存值
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set 设置缓存值
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonValue, expiration).Err()
}

// Delete 删除缓存
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists 检查缓存是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	cache map[string]cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[string]cacheItem),
	}
}

// Get 获取缓存值
func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	item, exists := m.cache[key]
	if !exists {
		return "", redis.Nil
	}

	if time.Now().After(item.expiration) {
		delete(m.cache, key)
		return "", redis.Nil
	}

	jsonValue, err := json.Marshal(item.value)
	if err != nil {
		return "", err
	}

	return string(jsonValue), nil
}

// Set 设置缓存值
func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.cache[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	return nil
}

// Delete 删除缓存
func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	delete(m.cache, key)
	return nil
}

// Exists 检查缓存是否存在
func (m *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.cache[key]
	return exists, nil
}
