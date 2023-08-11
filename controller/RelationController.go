package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/service"
	"github.com/goTouch/TicTok_SimpleVersion/util"
	"log"
	"net/http"
	"strconv"
)

// RelationAction 关注操作action_type : 1是关注，2是取消关注
func RelationAction(c *gin.Context) {
	userIdInt64, err := util.VerifyTokenReturnUserIdInt64(c)
	if err != nil {
		log.Println(err)
		return
	}

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

	err = service.Action(userIdInt64, toUserIdInt64, actionTypeInt)

}
