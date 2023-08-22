package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"tiktok_v2/pkg/viper"
	"time"
)

var (
	redisConfig = viper.Init("db")
	RedisClient = redis.Client{}
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Viper.GetString("redis.addr"), redisConfig.Viper.GetInt("redis.port")),
		Password:     redisConfig.Viper.GetString("redis.password"),
		DB:           redisConfig.Viper.GetInt("redis.db"),
		DialTimeout:  8 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 30 * time.Second,
		//配置连接池
		PoolSize:    10,
		PoolTimeout: 30 * time.Second,
	})
	//sync.Once 中的Do()传入一个函数作为参数，Do()会判断该函数是否已经被执行过，没有就执行
	//Do() 方法时会初始化 Redis 连接单例，后续的请求都会直接返回已经初始化好的连接对象。
	singletonRedis := sync.Once{}
	singletonRedis.Do(
		func() {
			RedisClient = *rdb
		},
	)
	return rdb
}

func init() {
	rdb := NewRedisClient()
	//连接成功了不需要返回值，连接失败了就要报错
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println("redis 连接成功！")

	//定期将数据更新到mysql

}
