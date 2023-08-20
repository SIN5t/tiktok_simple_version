package main

import (
	"log"
	"tiktok_v2/cmd/video/service"
	video "tiktok_v2/kitex_gen/video/videoservice"
)

func main() {
	svr := video.NewServer(new(service.VideoServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
