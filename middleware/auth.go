package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/service"
	"github.com/goTouch/TicTok_SimpleVersion/util"
)

// 强制鉴权，成功后会把userId写入context中，失败会直接终止请求
func AuthJWT(c *gin.Context) {
	tokenString := c.Query("token")
	userId, err := service.VerifyJWT(tokenString, util.JWTSecret())
	fmt.Println("userId", userId)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "用户未登录"})
		c.Abort()
		return
	}
	c.Set("userId", userId)
}

// 非强制鉴权，成功会将userId写入context中，失败会写入0
func AuthJWTOptional(c *gin.Context) {
	tokenString := c.Query("token")
	userId, err := service.VerifyJWT(tokenString, util.JWTSecret())
	fmt.Println("userId", userId)
	if err != nil {
		c.Set("userId", int64(0))
		return
	}
	c.Set("userId", userId)
}
