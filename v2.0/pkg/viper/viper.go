package viper

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Viper *viper.Viper
}

// Init 初始化配置文件，通过指定的名称，从指定路径中搜索出对应的yml配置文件 并返回带有这个yml参数的结构体
func Init(configName string) (config Config) {
	viperConfig := Config{Viper: new(viper.Viper)}
	viper := viperConfig.Viper
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")

	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	//读取配置文件中的内容
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return
	}
	return config
}
