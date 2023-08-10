package service

import (
	"context"
	"github.com/goTouch/TicTok_SimpleVersion/dao"
	"github.com/goTouch/TicTok_SimpleVersion/domain"
	"github.com/goTouch/TicTok_SimpleVersion/util"
	"log"
	"strconv"
	"time"
)

func FeedService(userIdInt64 int64, latestTimeInt64 int64) (videoList []domain.Video, nextTimeInt64 int64) {
	//将int64格式时间戳转为Time.time类型，以保证和数据库类型一致
	timeStamp := time.UnixMilli(userIdInt64)
	dao.DB.Model(&domain.Video{}).Preload("Author").
		Where("creat_time >= ?", timeStamp). // TODO 斟酌一下> 还是<
		Order("creat_time desc").            //该字段应该建一个索引提高效率
		Limit(3).                            //文档要求为30，这里设置小一点方便测试
		Find(&videoList)                     //保存到videoList中，最后返回给controller

	if len(videoList) == 0 {
		log.Println("FeedService查询数据库查到0条记录")
		return
	}

	// 返回这次视频最近的投稿时间-1，下次即可获取比这次视频旧的视频
	// videoList[len(videoList)-1]，因为排序过，这个就是最早的时间，将nextTime更新为当前最早视频时间减一，下个视频就从nextTime开始
	nextTimeInt64 = videoList[len(videoList)-1].CreatTime.UnixMilli() - 1

	//测试一下redis连接情况
	ping := dao.RedisClient.Ping(context.Background())
	log.Println(ping)

	for _, video := range videoList {
		// TODO 丰富Video的额外字段，例如author

		//查出每个视频对于当前用户的喜欢状态，已经视频作者的关注状态
		//注意前提是登入才能处理
		if userIdInt64 != -1 { //已登入
			//redis Set类型, key主要是当前访问用户的id，val是当前访问用户点赞的各个视频id
			// TODO 关注一下Hset保存以上信息的地方，下面是直接取出了
			isFavorite := dao.RedisClient.
				SIsMember(context.Background(), util.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), video.Id).
				Val()

			if isFavorite {
				//如果当前用户的点赞set中含有当前视频
				video.IsFavorite = true
			}
			//类似的，上面是点赞，这里是关注
			//key主要是当前用户的id，Hset的val是多个值：当前用户关注的作者的id

			isFollowed := dao.RedisClient.
				SIsMember(context.Background(), util.AuthorFollowedKeyPrefix+string(userIdInt64), video.AuthorId).
				Val()

			if isFollowed {
				//如果是关注的
				video.Author.IsFollow = true
			}
		}
		//未登入。默认显示没点赞、没关注。也就是直接数据库查询出来的结果：false

	}
	return
}
