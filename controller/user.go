package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/service"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	id, tokenString, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, domain.UserLoginResponse{
			Response: domain.Response{StatusCode: 0},
			UserId:   id,
			Token:    tokenString,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, tokenString, err := service.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, domain.UserLoginResponse{
			Response: domain.Response{StatusCode: 0},
			UserId:   id,
			Token:    tokenString,
		})
	}
}

func User(c *gin.Context) {
	id := c.Query("user_id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		id64 = 0
	}
	user, err := service.User(id64)
	if err != nil {
		c.JSON(http.StatusOK, domain.UserResponse{
			Response: domain.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, domain.UserResponse{
			Response: domain.Response{StatusCode: 0},
			User:     user,
		})
	}
}
