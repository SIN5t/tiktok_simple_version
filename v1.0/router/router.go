package router

import (
	"github.com/gin-gonic/gin"
	controller2 "github.com/goTouch/TicTok_SimpleVersion/v1.0/controller"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/middleware"
)

func InitRouter(r *gin.Engine) {
	// miniodata directory is used to serve static resources
	r.Static("/static", "./miniodata")

	apiR := r.Group("/douyin")

	// feed
	feedR := apiR.Group("/feed").Use(middleware.AuthJWTOptional)
	feedR.GET("/", controller2.Feed)

	// user
	userR := apiR.Group("/user")
	userR.POST("/register/", controller2.Register)
	userR.POST("/login/", controller2.Login)
	userR.GET("/", middleware.AuthJWTForce, controller2.User)

	// publish
	pubR := apiR.Group("/publish").Use(middleware.AuthJWTForce)
	pubR.POST("/action/", controller2.Publish)
	pubR.GET("/list", controller2.PublishList)

	// favorite
	favR := apiR.Group("/favorite").Use(middleware.AuthJWTForce)
	favR.POST("/action/", controller2.FavoriteAction)
	favR.GET("/list/", controller2.FavoriteList)

	// comment
	cmtR := apiR.Group("/comment").Use(middleware.AuthJWTForce)
	cmtR.POST("/action/", controller2.CommentAction)
	cmtR.GET("/list/", controller2.CommentList)

	// relation
	rltR := apiR.Group("/relation").Use(middleware.AuthJWTForce)
	rltR.POST("/action/", controller2.FollowAction)
	rltR.GET("/follow/list/", controller2.FollowList)
	rltR.GET("/follower/list/", controller2.FollowerList)
	rltR.GET("/friend/list/", controller2.FriendList)

	// message
	msgR := apiR.Group("/message").Use(middleware.AuthJWTForce)
	msgR.GET("/chat/", controller2.Chat)          // TODO
	msgR.POST("/action/", controller2.ChatAction) // TODO
}
