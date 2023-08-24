package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goForward/tictok_simple_version/domain"
	"github.com/goForward/tictok_simple_version/service"
)

func LoginLimit(c *gin.Context) {
	ipAddress := c.ClientIP()
	ok := service.LoginLimit(ipAddress)
	if !ok {
		c.JSON(http.StatusOK, domain.UserLoginResponse{
			Response: domain.Response{StatusCode: 1, StatusMsg: "操作过于频繁，请稍后再试"},
		})
		c.Abort()
	}
}
