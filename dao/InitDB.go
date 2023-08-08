package dao

import (
	"github.com/go-redis/redis/v8"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

/**
初始化数据库，包括redis和使用gorm
*/

var (
	DB          *gorm.DB
	RedisClient *redis.Client
	RdbToken    *redis.Client
)

const (
	numTokenDB = iota
)

func InitDB() {

	//datasource
	dsn := "root:123456@tcp(localhost:3306)/" +
		"tiktok?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local"
	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Println("InitDB中数据库初始化失败！")
		panic(err)
	}

	//创建数据库表格或更新已存在的表格
	err = DB.AutoMigrate(&domain.User{}, &domain.Video{}, &domain.Comment{})
	if err != nil {
		//return
		log.Println(err)
	}
	RdbToken = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numTokenDB,
	})
	// 创建 Redis 客户端配置
	redisConfig := &redis.Options{
		Addr:     "127.0.0.1:6379", // Redis 服务器地址和端口
		Password: "",               // Redis 认证密码，如果没有密码则为空字符串
		DB:       0,                // 选择使用的数据库，默认为 0
	}

	// 初始化 Redis 客户端
	RedisClient = redis.NewClient(redisConfig)

}
