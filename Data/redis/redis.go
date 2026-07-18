package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rdb *redis.Client
var ctx = context.Background()

func Init() error {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			viper.GetString("redis.host"),
			viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	if rdb == nil {
		return
	}
	err := rdb.Close()
	if err != nil {
		zap.L().Fatal("redis未能正常关闭", zap.Error(err))
	} else {
		zap.L().Info("redis正常退出")
	}
}
