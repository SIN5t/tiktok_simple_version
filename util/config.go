package util

import (
	"fmt"
	"os"
	"strings"
)

const (
	StaticRooterPrefix = "http://127.0.0.1:8080/" //

	AuthorBeLikedNum = "AUTHOR_BE_LIKED_NUM_KEY:"

	UserFollowHashPrefix    = "USER_FOLLOWS_KEY:"   //当前用户的关注hash列表，field是关注用户的id，value是对应的名字
	UserFollowersHashPrefix = "USER_FOLLOWERS_KEY:" //当前用户的粉丝列表。field是粉丝id,value是粉丝名字。

	VideoFavoriteKeyPrefix = "VIDEO_FAVORITE_KEY:" //用户角度，用处：该用户点赞视频
	VidioBucketName        = "vidio"
	PictureBucketName      = "picture"
)

const projectId = "tiktok:" // 项目标识符

// 获取项目环境变量，例如 projectEnv("mysql_user") 会获取 tiktok:mysql_user 这个环境变量的值
// 空值如果传入了defaultVal，会返回第一个defaultVal
func projectEnv(key string, defaultVal ...string) string {
	value := os.Getenv(projectId + key)
	if len(defaultVal) > 0 && value == "" {
		return defaultVal[0]
	}
	return value
}

var (
	mysql_user = projectEnv("mysql_user", "root")
	mysql_pswd = projectEnv("mysql_pswd", "123456")
	mysql_addr = projectEnv("mysql_addr", "localhost:3306")
	mysql_name = projectEnv("mysql_name", "tiktok")

	endpoint             = projectEnv("endpoint", "127.0.0.1:9000")
	accessKeyID          = projectEnv("accessKeyID", "minioadmin")
	secretAccessKey      = projectEnv("secretAccessKey", "minioadmin")
	redis_addr           = projectEnv("redis_addr", "localhost:6379", "192.168.157.128:6379")
	redis_addr_linux     = projectEnv("redis_addr", "192.168.157.128:6379")
	redis_pswd           = projectEnv("redis_pswd", "123456")
	redis_master_name    = projectEnv("redis_master_name", "mymaster")
	redis_sentinel_addrs = strings.Split(projectEnv("redi_sentinel_addrs", ":17000 :17001 :17002 "), "")

	jwt_secret = projectEnv("jwt_secret", "f05ad7412aa192ddc121ba50a64e585943b3e6d8fca4a3a19a8eea26e76496a7")
)

func GetMySQLDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local", mysql_user, mysql_pswd, mysql_addr, mysql_name)
	fmt.Println(dsn)
	return dsn
}

func GetRedisAddr() string {
	//return redis_addr
	return redis_addr_linux
}

func GetRedisPswd() string {
	return redis_pswd
}

func GetRedisMasterName() string {
	return redis_master_name
}

func GetRedisSentinelAddrs() []string {
	return redis_sentinel_addrs
}

func JWTSecret() string {
	return jwt_secret
}
func GetMinioEndpoint() string {
	return endpoint
}
func GetMinioAccessKeyID() string {
	return accessKeyID
}
func GetMinioSecretAccessKey() string {
	return secretAccessKey
}
