package controller

import (
	domain2 "github.com/goTouch/TicTok_SimpleVersion/v1.0/domain"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	userId := c.GetInt64("userId")
	//根据接口文档，前端传来的request中有token和latest_time， 这里一个用于存当前用户id，一个存下次视频时间戳
	//tokenReq := c.Query("token")
	latestTimeReq := c.Query("latest_time")                         //字符串类型
	latestTimeInt64, err := strconv.ParseInt(latestTimeReq, 10, 64) //转为时间戳
	if err != nil {
		c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: "时间戳格式错误"}) //定义1为错误的返回
		return
	}

	videoList, nextTimeInt64, err := service.FeedService(userId, latestTimeInt64)
	if err != nil {
		c.JSON(http.StatusOK, domain2.Response{
			StatusCode: 1, StatusMsg: err.Error(),
		})
		return
	}
	//说明查到了视频
	c.JSON(http.StatusOK, domain2.FeedResponse{
		Response:  domain2.Response{StatusCode: 0, StatusMsg: "刷新成功"},
		VideoList: videoList,
		NextTime:  nextTimeInt64,
	})

}
