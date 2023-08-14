package router

import (
	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/controller"
	"github.com/goTouch/TicTok_SimpleVersion/middleware"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiR := r.Group("/douyin")

	// feed
	feedR := apiR.Group("/feed").Use(middleware.AuthJWTOptional)
	feedR.GET("/", controller.Feed)

	// user
	userR := apiR.Group("/user")
	userR.POST("/register/", controller.Register)
	userR.POST("/login/", controller.Login)
	userR.GET("/", middleware.AuthJWT, controller.User)

	// publish
	pubR := apiR.Group("/publish").Use(middleware.AuthJWT)
	pubR.POST("/action/", controller.Publish)
	pubR.GET("/list") // TODO

	// favorite
	favR := apiR.Group("/favorite").Use(middleware.AuthJWT)
	favR.POST("/action/", controller.FavoriteAction)
	favR.GET("/list/", controller.FavoriteList)

	// comment
	cmtR := apiR.Group("/comment").Use(middleware.AuthJWT)
	cmtR.POST("/action/", controller.CommentAction)
	cmtR.GET("/list/", controller.CommentList)

	// relation
	rltR := apiR.Group("/relation").Use(middleware.AuthJWT)
	rltR.POST("/action/", controller.FollowAction)
	rltR.GET("/follow/list/", controller.FollowList)
	rltR.GET("/follower/list/", controller.FollowerList)
	rltR.GET("/friend/list/") // TODO

	// message
	msgR := apiR.Group("/message").Use(middleware.AuthJWT)
	msgR.GET("/chat")    // TODO
	msgR.POST("/action") // TODO
}
