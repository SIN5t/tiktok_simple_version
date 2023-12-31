package service

import (
	"context"
	"github.com/goForward/tictok_simple_version/config"
	"strconv"
	"time"

	"github.com/goForward/tictok_simple_version/dao"
	"github.com/goForward/tictok_simple_version/domain"
)

func FeedService(userIdInt64 int64, latestTimeInt64 int64) (videoList []domain.Video, nextTimeInt64 int64, err error) {

	userIdStr := strconv.FormatInt(userIdInt64, 10)
	//将int64格式时间戳转为Time.time类型，以保证和数据库类型一致
	timeStamp := time.UnixMilli(latestTimeInt64)
	dao.DB.Model(&domain.Video{}).Preload("Author").
		Where("creat_time <= ?", timeStamp). // 应该用小于
		Order("creat_time desc").            // 该字段应该建一个索引提高效率
		Limit(3).                            // 文档要求为30，这里设置小一点方便测试
		Find(&videoList)                     // 保存到videoList中，最后返回给controller

	if len(videoList) == 0 {
		nextTimeInt64 = time.Now().UnixMilli()
		return
	}

	// 返回这次视频最近的投稿时间-1，下次即可获取比这次视频旧的视频
	nextTimeInt64 = videoList[len(videoList)-1].CreatTime.UnixMilli() - 1
	url := dao.MinioClient.EndpointURL().String() + "/" + config.VideoBucketName + "/"
	picurl := dao.MinioClient.EndpointURL().String() + "/" + config.PictureBucketName + "/"
	for i := 0; i < len(videoList); i++ {

		video := &videoList[i]
		if videoList[i].CoverUrl == "" {
			videoList[i].CoverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
		} else {
			videoList[i].CoverUrl = picurl + videoList[i].CoverUrl
		}

		videoList[i].PlayUrl = url + videoList[i].PlayUrl
		//查出每个视频对于当前用户的喜欢状态，已经视频作者的关注状态
		//注意前提是登入才能处理
		if userIdInt64 != 0 { //已登入
			//是否点赞
			isFavorite := dao.RedisClient.
				SIsMember(context.Background(), config.VideoFavoriteKeyPrefix+userIdStr, video.Id).
				Val()

			if isFavorite {
				//如果当前用户的点赞set中含有当前视频
				video.IsFavorite = true
			}

			/*//视频被点赞总数
			dao.RedisClient.
				SCard(context.Background(), util.VideoFavoriteKeyPrefix+userIdStr)*/
			//关注
			isFollowed := dao.RedisClient.
				HExists(context.Background(), config.UserFollowHashPrefix+userIdStr, strconv.FormatInt(video.AuthorId, 10)).
				Val()

			if isFollowed {
				//如果当前作者是关注的作者
				video.Author.IsFollow = true
			}
		}

	}
	return
}
