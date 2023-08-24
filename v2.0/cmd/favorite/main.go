package main

import (
	"log"
	favorite2 "tiktok_v2/cmd/favorite/service"
	favorite "tiktok_v2/kitex_gen/favorite/favoriteservice"
)

func main() {
	svr := favorite.NewServer(new(favorite2.FavoriteServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
