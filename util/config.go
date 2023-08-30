package util

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const (
	AuthorBeLikedNum = "AUTHOR_BE_LIKED_NUM_KEY:"
	VideoBeLikedNum  = "VIDEO_BE_LIKED_NUM_KEY:"

	UserFollowHashPrefix    = "USER_FOLLOWS_KEY:"   // 当前用户的关注hash列表，field是关注用户的id，value是对应的名字
	UserFollowersHashPrefix = "USER_FOLLOWERS_KEY:" // 当前用户的粉丝列表。field是粉丝id,value是粉丝名字。

	VideoFavoriteKeyPrefix = "VIDEO_FAVORITE_KEY:" // 存用户点赞视频的ids

	UserMessageTimePrefix = "USER_MESSAGE_KEY:"

	TokenRefreshPrefix = "TOKEN_REFRESH:"

	VideoBucketName   = "video"
	PictureBucketName = "picture"
)

type Config struct {
	MysqlUser            string
	MysqlPswd            string
	MysqlAddr            string
	MysqlName            string
	RedisAddr            string
	RedisPswd            string
	RedisMasterName      string
	RedisSentinelAddrs   []string
	MinioEndpoint        string
	MinioAccessKeyID     string
	MinioSecretAccessKey string
	JWTSecret            string
}

var cfg Config

func InitConfig() {
	doc, err := os.ReadFile("./config.toml")
	if err != nil {
		panic(fmt.Errorf("读取配置文件 config.toml 失败: %w", err))
	}

	err = toml.Unmarshal(doc, &cfg)
	if err != nil {
		panic(fmt.Errorf("解析配置文件 config.toml 失败: %w", err))
	}
	fmt.Printf("%#v", cfg)
}

func GetMySQLDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local",
		cfg.MysqlUser, cfg.MysqlPswd, cfg.MysqlAddr, cfg.MysqlName)
	return dsn
}

func GetRedisAddr() string {
	return cfg.RedisAddr
}

func GetRedisPswd() string {
	return cfg.RedisPswd
}

func GetRedisMasterName() string {
	return cfg.RedisMasterName
}

func GetRedisSentinelAddrs() []string {
	return cfg.RedisSentinelAddrs
}

func JWTSecret() string {
	return cfg.JWTSecret
}

func GetMinioEndpoint() string {
	return cfg.MinioEndpoint
}

func GetMinioAccessKeyID() string {
	return cfg.MinioAccessKeyID
}

func GetMinioSecretAccessKey() string {
	return cfg.MinioSecretAccessKey
}

func Key(prefix string, val any) string {
	return fmt.Sprintf("%v%v", prefix, val)
}

// const projectId = "tiktok:" // 项目标识符

// // 获取项目环境变量，例如 projectEnv("mysql_user") 会获取 tiktok:mysql_user 这个环境变量的值
// // 空值如果传入了defaultVal，会返回第一个defaultVal
// func projectEnv(key string, defaultVal ...string) string {
// 	value := os.Getenv(projectId + key)
// 	if len(defaultVal) > 0 && value == "" {
// 		return defaultVal[0]
// 	}
// 	return value
// }

// var (
// 	mysql_user              = projectEnv("mysql_user", "root")
// 	mysql_pswd              = projectEnv("mysql_pswd", "123456")
// 	mysql_addr              = projectEnv("mysql_addr", "localhost:3306")
// 	mysql_name              = projectEnv("mysql_name", "tiktok")
// 	minio_endpoint          = projectEnv("minio_endpoint", "127.0.0.1:9000")
// 	minio_access_key_id     = projectEnv("minio_access_key_id", "minioadmin")
// 	minio_secret_access_key = projectEnv("minio_secret_access_key", "minioadmin")
// 	redis_addr              = projectEnv("redis_addr", "localhost:6379")
// 	redis_pswd              = projectEnv("redis_pswd", "123456")
// 	redis_master_name       = projectEnv("redis_master_name", "mymaster")
// 	redis_sentinel_addrs    = strings.Split(projectEnv("redi_sentinel_addrs", ":17000 :17001 :17002 "), "")
// 	jwt_secret              = projectEnv("jwt_secret", "f05ad7412aa192ddc121ba50a64e585943b3e6d8fca4a3a19a8eea26e76496a7")
// )
