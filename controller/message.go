package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/service"
	"log"
	"net/http"
	"strconv"
)

func Chat(c *gin.Context) {
	fromUserId := c.GetInt64("userId")
	toUserId, err1 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "toUserId解析失败"})
		return
	}
	list, err := service.ChatList(fromUserId, toUserId)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "获取聊天记录失败"})
		return
	}
	c.JSON(http.StatusOK, domain.ChatResponse{
		Response:    domain.Response{StatusCode: 0},
		MessageList: list,
	})
	return
}
func ChatAction(c *gin.Context) {
	fromUserId := c.GetInt64("userId")
	toUserId, err1 := strconv.ParseInt(c.PostForm("to_user_id"), 10, 64)
	content := c.PostForm("content")
	//actionType := c.PostForm("action_type")
	actionType, err2 := strconv.ParseInt(c.PostForm("action_type"), 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "userId或actiontype解析失败"})
		return
	}
	if content == "" || actionType != 1 {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "参数错误"})
		return
	}
	message, err := service.AddMessage(fromUserId, toUserId, content)
	if err != nil {
		log.Println(message, err)
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: message + err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.Response{StatusCode: 0, StatusMsg: message})
	return
}
