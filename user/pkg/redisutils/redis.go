package redisutils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"user/config"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func InitRedisConnect() (err error) {
	redisConfig := config.GetConfig().Redis
	var (
		DefaultPoolMaxIdle         = 100
		DefaultPoolSize            = 500
		DefaultPoolMaxConnLifetime = 10 // 单位s
		DefaultWriteTimeout        = 5  // 写入超时时间
		DefaultReadTimeout         = 2  // 读取超时时间
	)
	// 拼接一下redis client要用的地址
	addr := fmt.Sprintf("%s:%d", redisConfig.Client.Host, redisConfig.Client.Port)
	// name := config.Server.Name
	rdbOpts := &redis.Options{
		Addr:     addr,
		Password: redisConfig.Client.Password,
		DB:       redisConfig.Client.DB,
	}

	// 最大空闲链接
	maxIdle := DefaultPoolMaxIdle
	if redisConfig.Pool.MaxIdle > 0 {
		maxIdle = redisConfig.Pool.MaxIdle
	}
	rdbOpts.MaxIdleConns = maxIdle

	// 池子大小
	poolSize := DefaultPoolSize
	if redisConfig.Pool.PoolSize > 0 {
		poolSize = redisConfig.Pool.PoolSize
	}
	rdbOpts.PoolSize = poolSize

	// 最大链接时间
	maxConnectLifeTime := time.Second * time.Duration(DefaultPoolMaxConnLifetime)
	if redisConfig.Pool.MaxConnLifetime > 0 {
		maxConnectLifeTime = time.Second * time.Duration(redisConfig.Pool.MaxConnLifetime)
	}
	rdbOpts.ConnMaxLifetime = maxConnectLifeTime

	// read/write timeout
	readTimeout := time.Second * time.Duration(DefaultReadTimeout)
	if redisConfig.Client.ReadTimeout > 0 {
		readTimeout = time.Second * time.Duration(redisConfig.Client.ReadTimeout)
	}
	rdbOpts.ReadTimeout = readTimeout

	writeTimeout := time.Second * time.Duration(DefaultWriteTimeout)
	if redisConfig.Client.WriteTimeout > 0 {
		writeTimeout = time.Second * time.Duration(redisConfig.Client.WriteTimeout)
	}
	rdbOpts.WriteTimeout = writeTimeout

	// 获取到db
	db := redis.NewClient(rdbOpts)
	client = db

	return nil
}

// Set 设置键值对，并可设置过期时间
// 参数：
// - ctx: 上下文
// - key: 缓存键
// - value: 缓存值
// - expiration: 过期时间（0 表示不过期）
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

// Get 获取键对应的值
// 返回：
// - string: 值（键不存在时为空字符串）
// - bool: 键是否存在
// - error: 操作错误（网络错误等）
func Get(ctx context.Context, key string) (string, bool, error) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil // 键不存在
		}
		return "", false, fmt.Errorf("获取缓存失败: %w", err)
	}
	return val, true, nil
}

// SetNX 设置键值对（仅当键不存在时）
// 返回：
// - bool: 是否设置成功
// - error: 操作错误
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return client.SetNX(ctx, key, value, expiration).Result()
}

// Delete 删除键
func Delete(ctx context.Context, key string) error {
	return client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	count, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查缓存存在失败: %w", err)
	}
	return count > 0, nil
}

// Expire 设置键的过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return client.Expire(ctx, key, expiration).Err()
}

// Incr 对键的值进行递增
func Incr(ctx context.Context, key string) (int64, error) {
	return client.Incr(ctx, key).Result()
}

// Decr 对键的值进行递减
func Decr(ctx context.Context, key string) (int64, error) {
	return client.Decr(ctx, key).Result()
}

func Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return client.Eval(ctx, script, keys, args...).Result()
}
