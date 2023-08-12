package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goTouch/TicTok_SimpleVersion/dao"
	"github.com/goTouch/TicTok_SimpleVersion/router"
)

func main() {
	dao.InitMinio()
	dao.InitDB()
	fmt.Print("数据库执行成功")

	//go service.RunMessageServer()

	r := gin.Default()

	router.InitRouter(r)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
