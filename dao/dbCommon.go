package dao

import (
	"github.com/bwmarrin/snowflake"
)

var (
	//RedisSyncNode, _ = snowflake.NewNode(1)
	UserNode, _  = snowflake.NewNode(2)
	VideoNode, _ = snowflake.NewNode(3)
)
