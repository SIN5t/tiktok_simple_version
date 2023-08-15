package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	dao "github.com/goTouch/TicTok_SimpleVersion/dao"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/service"
	"github.com/goTouch/TicTok_SimpleVersion/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Chat(c *gin.Context) {

	fromUserId := c.GetInt64("userId")
	//msgTime := c.GetInt64("pre_msg_time")
	msgTime, err2 := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)

	toUserId, err1 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "toUserId解析失败"})
		return
	}
	//log.Println(time.Now().UnixMilli())
	//log.Println(c.Query("pre_msg_time"))
	//
	//log.Println(msgTime < time.Now().UnixMilli())
	queryTime := time.Now().UnixMilli()

	if msgTime > queryTime {
		time, _ := dao.RedisClient.Get(c, util.UserMessageTimePrefix+strconv.FormatInt(fromUserId, 10)+":"+strconv.FormatInt(toUserId, 10)).Result()
		msgTime, _ = strconv.ParseInt(time, 10, 64)
	}

	list, err := service.ChatList(fromUserId, toUserId, msgTime)
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
	fmt.Println(time.Now().Format(time.RFC3339))
	fromUserId := c.GetInt64("userId")
	toUserId, err1 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	content := c.Query("content")
	//actionType := c.PostForm("action_type")
	actionType, err2 := strconv.ParseInt(c.Query("action_type"), 10, 64)
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
	queryTime := time.Now().UnixMilli()
	dao.RedisClient.Set(c, util.UserMessageTimePrefix+strconv.FormatInt(fromUserId, 10)+":"+strconv.FormatInt(toUserId, 10), strconv.FormatInt(queryTime, 10), 0)

	c.JSON(http.StatusOK, domain.Response{StatusCode: 0, StatusMsg: message})
	return
}
