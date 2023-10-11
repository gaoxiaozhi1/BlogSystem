package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
		PoolSize: 100,      // 连接池大小
	})
	// 判断是否超时，是否连接成功
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		logrus.Error(err)
		return
	}
}

func main() {
	// 创建存字符串
	err := rdb.Set("xxx1", "value1", 10*time.Second).Err()
	fmt.Println(err)
	// 查询所有
	cmd := rdb.Keys("*")
	keys, err := cmd.Result()
	fmt.Println(keys, err)

}
