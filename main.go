package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/goForward/tictok_simple_version/dao"
	"github.com/goForward/tictok_simple_version/router"
	"github.com/goForward/tictok_simple_version/util"
)

func main() {
	util.InitConfig()
	log.Println("配置读取成功")

	dao.InitMinio()
	dao.InitDB()
	log.Println("数据库执行成功")

	//go service.RunMessageServer()

	r := gin.Default()

	router.InitRouter(r)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
