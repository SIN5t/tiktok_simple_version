package router

import (
	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/controller"
	"github.com/goTouch/TicTok_SimpleVersion/middleware"
)

func InitRouter(r *gin.Engine) {
	// miniodata directory is used to serve static resources
	r.Static("/static", "./miniodata")

	apiR := r.Group("/douyin")

	// feed
	feedR := apiR.Group("/feed").Use(middleware.AuthJWTOptional)
	feedR.GET("/", controller.Feed)

	// user
	userR := apiR.Group("/user")
	userR.POST("/register/", controller.Register)
	userR.POST("/login/", controller.Login)
	userR.GET("/", middleware.AuthJWTForce, controller.User)

	// publish
	pubR := apiR.Group("/publish").Use(middleware.AuthJWTForce)
	pubR.POST("/action/", controller.Publish)
	pubR.GET("/list")

	// favorite
	favR := apiR.Group("/favorite").Use(middleware.AuthJWTForce)
	favR.POST("/action/", controller.FavoriteAction)
	favR.GET("/list/", controller.FavoriteList)

	// comment
	cmtR := apiR.Group("/comment").Use(middleware.AuthJWTForce)
	cmtR.POST("/action/", controller.CommentAction)
	cmtR.GET("/list/", controller.CommentList)

	// relation
	rltR := apiR.Group("/relation").Use(middleware.AuthJWTForce)
	rltR.POST("/action/", controller.FollowAction)
	rltR.GET("/follow/list/", controller.FollowList)
	rltR.GET("/follower/list/", controller.FollowerList)
	rltR.GET("/friend/list/", controller.FriendList)

	// message
	msgR := apiR.Group("/message").Use(middleware.AuthJWTForce)
	msgR.GET("/chat/", controller.Chat)          // TODO
	msgR.POST("/action/", controller.ChatAction) // TODO
}
