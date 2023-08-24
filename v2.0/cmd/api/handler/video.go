package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok_v2/cmd/api/rpc"
	"tiktok_v2/internal/response"
	kitexVideo "tiktok_v2/kitex_gen/video"
	"time"
)

func Feed(c *gin.Context) {
	//handler层获取前端提交的数据
	token := c.Query("token")
	lastTimeStr := c.Query("latest_time")

	timeStamp := time.Now().UnixMilli() //考虑lastTimeStr为空

	if lastTimeStr != "" {
		//使用客户的时间戳
		timeStamp, _ = strconv.ParseInt(lastTimeStr, 10, 64)
	}

	//构建请求体 使用地址
	req := &kitexVideo.FeedRequest{
		LatestTime: timeStamp,
		Token:      token,
	}
	//rpc调用，构建一个与微服务service层远程联系的client

	res, err := rpc.FeedClient(context.Background(), req)
	if err != nil {
		// TODO 日志
		c.JSON(http.StatusOK, response.Base{
			StatusCode: 1,
			StatusMsg:  "服务器连接失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Feed{
		Base:      response.Base{StatusCode: 0, StatusMsg: "刷新成功"},
		NextTime:  res.NextTime,
		VideoList: res.VideoList,
	})

}
