package main

import (
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
	"tiktok_v2/cmd/video/service"
	"tiktok_v2/kitex_gen/video/videoservice"
	"tiktok_v2/pkg/etcd"
	"tiktok_v2/pkg/viper"
)

var (
	config      = viper.Init("video")
	serviceName = config.Viper.GetString("server.name") //server.WithServerBasicInfo用到
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	//signingKey  = config.Viper.GetString("JWT.signingKey")

)

func main() {

	//服务注册
	registry, err2 := etcd.NewEtcdRegistry([]string{etcdAddr})
	if err2 != nil {
		log.Fatal(err2.Error())
	}
	//tcp的resolver，返回一个tcpAddr,含有ip+port
	addr, err2 := net.ResolveTCPAddr("tcp", serviceAddr)
	if err2 != nil {
		log.Fatal(err2.Error())
	}
	videoserver := videoservice.NewServer(
		new(service.VideoServiceImpl), //这个service就是mvc中的service
		server.WithServiceAddr(addr),  //tcp 的地址
		server.WithRegistry(registry), //注册中心
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)

	if err := videoserver.Run(); err != nil {
		logger.Fatalf("%v stopped with error: %v", serviceName, err.Error())
	}

}
