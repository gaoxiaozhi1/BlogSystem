package core

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gvb_server/global"
	"time"
)

func ConnectRedis() *redis.Client {
	return ConnectRedisDB(0)
}

// ConnectRedisDB 根据不同的DB连接不同的redis
func ConnectRedisDB(db int) *redis.Client {
	redisConf := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr(),
		Password: redisConf.Password, // no password set
		DB:       db,                 // use default DB
		PoolSize: redisConf.PoolSize, // 连接池大小
	})

	// 判断是否超时，是否连接成功
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		logrus.Error("redis连接失败 %s", redisConf.Addr())
		return nil
	}
	return rdb
}
