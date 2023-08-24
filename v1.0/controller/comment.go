package controller

import (
	domain2 "github.com/goTouch/TicTok_SimpleVersion/v1.0/domain"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CommentAction(c *gin.Context) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: "非法视频id"})
		return
	}

	userId := c.GetInt64("userId")
	actionType := c.Query("action_type")

	if actionType == "1" {
		commentText := c.Query("comment_text")
		comment, err := service.AddComment(videoId, userId, commentText)
		if err != nil {
			c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: "创建评论失败"})
		} else {
			c.JSON(http.StatusOK, domain2.CommentResponse{
				Response: domain2.Response{StatusCode: 0},
				Comment:  comment,
			})
		}
	} else if actionType == "2" {
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: "非法评论id"})
			return
		}

		err = service.DeleteComment(commentId)
		if err != nil {
			c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: "删除评论失败"})
		} else {
			c.JSON(http.StatusOK, domain2.Response{StatusCode: 0})
		}
	} else {
		c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: "无效操作"})
	}
}

func CommentList(c *gin.Context) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: "非法视频id"})
		return
	}

	commentList, err := service.CommentList(videoId)
	if err != nil {
		c.JSON(http.StatusOK, domain2.Response{StatusCode: 1, StatusMsg: "非法视频id"})
	} else {
		c.JSON(http.StatusOK, domain2.CommentListResponse{
			Response:    domain2.Response{StatusCode: 0},
			CommentList: commentList,
		})
	}
}
