package main

import (
	"log"
	relation2 "tiktok_v2/cmd/relation/service"
	relation "tiktok_v2/kitex_gen/relation/relationservice"
)

func main() {
	svr := relation.NewServer(new(relation2.RelationServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
