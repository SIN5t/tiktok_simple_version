package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goForward/tictok_simple_version/domain"
	"github.com/goForward/tictok_simple_version/service"
	"github.com/goForward/tictok_simple_version/util"
)

func authJWT(c *gin.Context, force bool) {
	// 获取token
	tokenString := c.PostForm("token")
	if tokenString == "" {
		tokenString = c.Query("token")
	}

	// 鉴权
	userId, err := service.VerifyJWT(tokenString, util.JWTSecret())
	fmt.Println("userId", userId)
	if err != nil { // 未登录
		if force { // 强制
			c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "用户未登录"})
			c.Abort()
		} // 非强制不处理
	} else { // 已登录
		// 刷新token
		if err = service.RefreshJWT(userId); err != nil {
			c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: err.Error()})
			c.Abort()
		}
	}

	// 用户id写入context中
	c.Set("userId", userId)
}

// 强制鉴权，成功后会把userId写入context中，失败会直接终止请求
func AuthJWTForce(c *gin.Context) {
	authJWT(c, true)
}

// 非强制鉴权，成功会将userId写入context中，失败会写入0
func AuthJWTOptional(c *gin.Context) {
	authJWT(c, false)
}
