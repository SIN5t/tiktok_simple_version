package main

import (
	dao2 "github.com/goTouch/TicTok_SimpleVersion/v1.0/dao"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/router"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/util"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	util.InitConfig()
	log.Println("配置读取成功")

	dao2.InitMinio()
	dao2.InitDB()
	log.Println("数据库执行成功")

	//go service.RunMessageServer()

	r := gin.Default()

	router.InitRouter(r)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
