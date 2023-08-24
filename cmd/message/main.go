package main

import (
	"log"
	message2 "tiktok_v2/cmd/message/service"
	message "tiktok_v2/kitex_gen/message/messageservice"
)

func main() {
	svr := message.NewServer(new(message2.MessageServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
