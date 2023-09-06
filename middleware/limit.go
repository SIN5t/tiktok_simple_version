package middleware

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/goForward/tictok_simple_version/dao"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goForward/tictok_simple_version/domain"
)

func OpsLimit(c *gin.Context) {
	ipAddress := c.ClientIP()
	times, err := dao.RedisClient.Get(context.Background(), ipAddress).Int64()
	if err != nil && err != redis.Nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if times > 10 {
		c.JSON(http.StatusOK, domain.UserLoginResponse{
			Response: domain.Response{StatusCode: 1, StatusMsg: "10s内操作过于频繁，请稍后再试"},
		})
		c.Abort()
		return
	}

	err = dao.RedisClient.Set(context.Background(), ipAddress, times+1, time.Second*10).Err()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Next()
}
