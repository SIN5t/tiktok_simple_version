package controller

import (
	domain2 "github.com/goTouch/TicTok_SimpleVersion/v1.0/domain"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	id, tokenString, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, domain2.UserLoginResponse{
			Response: domain2.Response{StatusCode: 0},
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
		c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, domain2.UserLoginResponse{
			Response: domain2.Response{StatusCode: 0},
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
		c.JSON(http.StatusOK, domain2.UserResponse{
			Response: domain2.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, domain2.UserResponse{
			Response: domain2.Response{StatusCode: 0},
			User:     user,
		})
	}
}
