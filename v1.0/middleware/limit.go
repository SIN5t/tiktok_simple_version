package middleware

import (
	domain2 "github.com/goTouch/TicTok_SimpleVersion/v1.0/domain"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginLimit(c *gin.Context) {
	ipAddress := c.ClientIP()
	ok := service.LoginLimit(ipAddress)
	if !ok {
		c.JSON(http.StatusOK, domain2.UserLoginResponse{
			Response: domain2.Response{StatusCode: 1, StatusMsg: "操作过于频繁，请稍后再试"},
		})
		c.Abort()
	}
}
