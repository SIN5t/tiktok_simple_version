package redis

import (
	"github.com/go-co-op/gocron"
	"time"
)

//数据库与redis一致性同步处理

func SyncFavoriteToMysql() {
	scheduler := gocron.NewScheduler(time.Local) ///定时任务
	scheduler.Every(10).
		Tag("favoriteRedis").
		Seconds().
		Do(SycFavoriteToMysql)
	scheduler.StartAsync() //异步执行
}

func SycFavoriteToMysql() {
	//1. 视频角度： 被点赞个数，先于是否使用redis

	//2.  用户角度： 该用户点赞的视频id、该用户点赞的总数（前端没做）

}
