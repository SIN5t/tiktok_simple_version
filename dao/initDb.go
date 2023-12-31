package dao

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/goForward/tictok_simple_version/config"
	"github.com/goForward/tictok_simple_version/domain"
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

	//RedisSyncNode, _ = snowflake.NewNode(1)
	UserNode, _  = snowflake.NewNode(2)
	VideoNode, _ = snowflake.NewNode(3)
)

func InitDB() {

	//datasource
	dsn := config.GetMySQLDSN()

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
	err = DB.AutoMigrate(&domain.User{}, &domain.Video{}, &domain.Comment{}, &domain.Message{}, &domain.UserRedisSync{})
	if err != nil {
		//return
		log.Println(err)
	}

	// *********************   redis   ************************************************

	// 创建 Redis 客户端配置
	redisConfig := &redis.Options{
		Addr:     config.GetRedisAddr(), // Redis 服务器地址和端口
		Password: config.GetRedisPswd(), // Redis 认证密码，如果没有密码则为空字符串
		DB:       0,                     // 选择使用的数据库，默认为 0
	}

	// 初始化 Redis 客户端
	RedisClient = redis.NewClient(redisConfig)
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Fatal("redis连接失败" + err.Error())
	}
	log.Println("successfully connected to Redis server!")

	//开启定时同步到数据库
	if err = ScheduleSyncFavVideoList(); err != nil {
		log.Println(err.Error())
	}
	if err = ScheduleSyncRelation(); err != nil {
		log.Println(err.Error())
	}
	if err = ScheduleSyncVideoBeLikedNum(); err != nil {
		log.Println(err.Error())
	}
	if err = ScheduleSyncAuthorBeLikedNum(); err != nil {
		log.Println(err.Error())
	}

	log.Println("MySQL synchronization is enabled.")
}
