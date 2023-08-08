package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/service"
	"net/http"
)

type UserLoginResponse struct {
	domain.Response
	Token int64 `json:"token"`
}

type UserResponse struct {
	domain.Response
	User domain.User `json:"user"`
}

// LoginLimit 中间件，限制注册登录操作过于频繁。
func LoginLimit(c *gin.Context) {
	ipAddress := c.ClientIP()
	ok := service.LoginLimit(ipAddress)
	if !ok {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: domain.Response{StatusCode: 1, StatusMsg: "操作过于频繁，请稍后再试"},
		})
		c.Abort()
	}
}
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	id, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		token := service.SaveTokenToRedis(id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: domain.Response{StatusCode: 0},
			Token:    token,
		})
	}
}
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, err := service.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		token := service.SaveTokenToRedis(id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: domain.Response{StatusCode: 0},
			Token:    token,
		})
	}
}
