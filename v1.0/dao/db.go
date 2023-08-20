package dao

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	domain2 "github.com/goTouch/TicTok_SimpleVersion/v1.0/domain"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
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
	dsn := util.GetMySQLDSN()

	//配置日志
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      ormLogger, //日志配置
		PrepareStmt: true,
	})
	if err != nil {
		log.Println("InitDB中数据库初始化失败！")
		panic(err)
	}

	//设置数据库参数
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(20)  //空闲时连接池
	sqlDB.SetMaxOpenConns(100) //最大打开数
	sqlDB.SetConnMaxLifetime(30 * time.Second)

	//创建数据库表格或更新已存在的表格
	err = DB.AutoMigrate(&domain2.User{}, &domain2.Video{}, &domain2.Comment{}, &domain2.Message{})
	if err != nil {
		//return
		log.Println(err)
	}
	RdbToken = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    util.GetRedisMasterName(),
		SentinelAddrs: util.GetRedisSentinelAddrs(),
		DB:            numTokenDB,
	})
	// 创建 Redis 客户端配置
	redisConfig := &redis.Options{
		Addr:     util.GetRedisAddr(), // Redis 服务器地址和端口
		Password: util.GetRedisPswd(), // Redis 认证密码，如果没有密码则为空字符串
		DB:       0,                   // 选择使用的数据库，默认为 0
	}

	// 初始化 Redis 客户端
	RedisClient = redis.NewClient(redisConfig)

}
