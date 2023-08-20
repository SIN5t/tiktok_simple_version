package main

import (
	"log"
	user2 "tiktok_v2/cmd/user/service"
	user "tiktok_v2/kitex_gen/user/userservice"
)

func main() {
	svr := user.NewServer(new(user2.UserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
