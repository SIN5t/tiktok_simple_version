package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/service"
)

//Relation

// FollowAction 关注操作action_type : 1是关注，2是取消关注
func FollowAction(c *gin.Context) {
	userIdInt64 := c.GetInt64("userId")

	toUserIdStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")

	toUserIdInt64, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		log.Println("err parsing userId")
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  "err parsing userId",
		})
	}
	//忽略错误，前端传来的关注按钮一般不会错
	actionTypeInt, _ := strconv.Atoi(actionTypeStr)

	err = service.FollowAction(userIdInt64, toUserIdInt64, actionTypeInt)
	if err != nil {
		//出错了要提示前端app
		c.JSON(http.StatusOK, domain.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//没出错也要显示关注成功
	c.JSON(http.StatusOK, domain.Response{
		StatusCode: 0,
		StatusMsg:  "成功关注！",
	})
}

// FollowList 查询用户关注列表
func FollowList(c *gin.Context) {
	//TODO 思考，本来所有接口都应该校验一下token，并刷新token过期时间，是否有全局中间件的做法？
	userIdStr := c.Query("user_id")
	//错误处理
	if userIdStr == "" {
		//TODO 可以放频率统计中间件
		log.Println("userId = nil")
		c.JSON(http.StatusOK, domain.UserFollowListResponse{
			Response: domain.Response{
				StatusCode: 1,
				StatusMsg:  "用户id不存在！",
			},
			UserFollowList: nil,
		})
		return
	}

	//id正确，开始查询
	userList, err := service.FollowList(userIdStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.UserFollowListResponse{
			Response:       domain.Response{StatusCode: 1, StatusMsg: err.Error()},
			UserFollowList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, domain.UserFollowListResponse{
		Response:       domain.Response{StatusCode: 0, StatusMsg: "成功查询到用户的关注列表"},
		UserFollowList: userList,
	})
}

// FollowerList 查询用户的粉丝列表
func FollowerList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	if userIdStr == "" {
		c.JSON(http.StatusOK, domain.UserFollowListResponse{
			Response:       domain.Response{StatusCode: 1, StatusMsg: "wrong user id"},
			UserFollowList: nil,
		})
	}
	userList, err := service.FollowerList(userIdStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, domain.UserFollowListResponse{
			Response:       domain.Response{StatusCode: 1, StatusMsg: err.Error()},
			UserFollowList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, domain.UserFollowListResponse{
		Response:       domain.Response{StatusCode: 0, StatusMsg: "成功查询到用户的粉丝列表"},
		UserFollowList: userList,
	})
}
