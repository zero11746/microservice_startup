package initialize

import (
	"common/applog"
	"common/env"
	"context"
	"log"
	"user/config"
	"user/pkg/database"
	"user/pkg/mongodbutils"
	"user/pkg/redisutils"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

// MustInit 组件初始化，若失败，会panic
func MustInit(ctx context.Context) {
	err := env.InitEnvConfig()
	if err != nil {
		panic(err)
		return
	}

	err = config.InitConfig("config.toml")
	if err != nil {
		panic(err)
		return
	}

	err = applog.InitLoggers(config.GetConfig().AppLog)
	if err != nil {
		panic(err)
		return
	}

	err = redisutils.InitRedisConnect()
	if err != nil {
		panic(err)
		return
	}

	err = database.InitMysqlConnect()
	if err != nil {
		panic(err)
		return
	}

	err = mongodbutils.InitMongoConnect(ctx)
	if err != nil {
		panic(err)
		return
	}
}
