package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goForward/tictok_simple_version/domain"
	"github.com/goForward/tictok_simple_version/service"
)

// FavoriteAction
func FavoriteAction(c *gin.Context) {

	//验证token，合法的话返回userId
	userIdInt64 := c.GetInt64("userId")

	videoIdInt64, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: "获取视频id失败！"})
		log.Println("出现无法解析成64位整数的视频id")
		return
	}

	//actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32) //这个函数返回的就是64位的
	actionType, err := strconv.Atoi(c.Query("action_type"))

	if err != nil {
		return
	}

	err = service.Favorite(videoIdInt64, userIdInt64, int32(actionType))
	if err != nil {
		c.JSON(http.StatusOK, domain.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.Response{StatusCode: 0, StatusMsg: "点赞成功"})

}

// FavoriteList 登录用户的所有点赞视频。
func FavoriteList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userIdInt64, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, domain.VideoListResponse{
			Response:  domain.Response{StatusCode: 1, StatusMsg: err.Error()},
			VideoList: nil,
		})
	}
	videoList, err := service.FavoriteList(userIdInt64)
	if err != nil {
		c.JSON(http.StatusOK, domain.VideoListResponse{
			Response:  domain.Response{StatusCode: 1, StatusMsg: "视频错误"},
			VideoList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, domain.VideoListResponse{
		Response:  domain.Response{StatusCode: 0},
		VideoList: videoList,
	})
}
