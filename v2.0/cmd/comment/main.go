package main

import (
	"log"
	"tiktok_v2/cmd/comment/service"
	comment "tiktok_v2/kitex_gen/comment/commentservice"
)

func main() {
	svr := comment.NewServer(new(service.CommentServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
