package mongodbutils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"user/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client    *mongo.Client
	defaultDB string // 存储配置文件中的默认数据库名
)

// InitMongoConnect 初始化MongoDB连接
func InitMongoConnect(ctx context.Context) (err error) {
	mongoConfig := config.GetConfig().Mongo
	// 默认配置值
	var (
		defaultConnectTimeout  = 5 // 秒
		defaultMaxPoolSize     = uint64(100)
		defaultMinPoolSize     = uint64(10)
		defaultMaxConnIdleTime = 300 // 秒
	)

	// 校验默认数据库名（必须配置）
	if mongoConfig.Database == "" {
		return errors.New("mongodb config: Database field is required")
	}
	defaultDB = mongoConfig.Database // 存储默认数据库名

	// 填充默认值（配置未设置时）
	if mongoConfig.ConnectTimeout == 0 {
		mongoConfig.ConnectTimeout = defaultConnectTimeout
	}
	if mongoConfig.MaxPoolSize == 0 {
		mongoConfig.MaxPoolSize = defaultMaxPoolSize
	}
	if mongoConfig.MinPoolSize == 0 {
		mongoConfig.MinPoolSize = defaultMinPoolSize
	}
	if mongoConfig.MaxConnIdleTime == 0 {
		mongoConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	}

	// 构建连接URI
	uri := fmt.Sprintf("mongodb://%s:%d", mongoConfig.Host, mongoConfig.Port)
	clientOpts := options.Client().ApplyURI(uri)

	// 设置认证信息
	if mongoConfig.Username != "" && mongoConfig.Password != "" {
		clientOpts.SetAuth(options.Credential{
			Username:   mongoConfig.Username,
			Password:   mongoConfig.Password,
			AuthSource: mongoConfig.AuthSource, // 默认为admin，配置中可自定义
		})
	}

	// 设置连接池参数
	clientOpts.SetMaxPoolSize(mongoConfig.MaxPoolSize)
	clientOpts.SetMinPoolSize(mongoConfig.MinPoolSize)
	clientOpts.SetMaxConnIdleTime(time.Duration(mongoConfig.MaxConnIdleTime) * time.Second)
	clientOpts.SetConnectTimeout(time.Duration(mongoConfig.ConnectTimeout) * time.Second)

	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("mongodb connect failed: %w", err)
	}

	// 验证连接
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return fmt.Errorf("mongodb ping failed: %w", err)
	}

	client.StartSession()

	return nil
}

// GetDatabase 获取数据库实例（默认使用配置中的数据库，支持传入自定义库名）
func GetDatabase(customDB ...string) *mongo.Database {
	if len(customDB) > 0 && customDB[0] != "" {
		return client.Database(customDB[0])
	}
	return client.Database(defaultDB) // 默认使用配置的数据库
}

// GetCollection 获取集合实例（默认使用配置的数据库，支持指定自定义库名）
func GetCollection(collectionName string, customDB ...string) *mongo.Collection {
	return GetDatabase(customDB...).Collection(collectionName)
}
