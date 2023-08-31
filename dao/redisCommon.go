package dao

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/goForward/tictok_simple_version/config"
	"github.com/goForward/tictok_simple_version/domain"
	"log"
	"strconv"
	"strings"
	"time"
)

//多协程异步处理 数据库与redis一致性

func ScheduleSyncFavVideoList() error {
	scheduler := gocron.NewScheduler(time.Local) ///定时任务
	_, err := scheduler.Every(60).
		Tag("favoriteRedis").
		Seconds().
		Do(SyncFavVideoList)
	if err != nil {
		return err
	}
	scheduler.StartAsync() //异步执行
	return nil
}
func ScheduleSyncRelation() error {
	scheduler := gocron.NewScheduler(time.Local)
	_, err := scheduler.Every(60).
		Tag("relationRedis").
		Second().
		Do(SyncRelationToMysql)
	if err != nil {
		return err
	}
	scheduler.StartAsync()
	return nil
}
func ScheduleSyncVideoBeLikedNum() error {
	scheduler := gocron.NewScheduler(time.Local)
	_, err := scheduler.Every(10).
		Tag("videoLikedNumRedis").
		Second().
		Do(SyncVideoBeLikedNum)
	if err != nil {
		return err
	}
	scheduler.StartAsync()
	return nil
}
func ScheduleSyncAuthorBeLikedNum() error {
	scheduler := gocron.NewScheduler(time.Local)
	_, err := scheduler.Every(10).
		Tag("authorLikedNumRedis").
		Second().
		Do(SyncAuthorLikedNum)
	if err != nil {
		return err
	}
	scheduler.StartAsync()
	return nil
}

// **********************具体逻辑实现********************************//

// SyncFavVideoList 点赞视频列表同步、用户点赞个数同步
func SyncFavVideoList() (err error) {
	//点赞视频列表同步
	matchPattern := config.VideoFavoriteKeyPrefix + "*"
	cursor := uint64(0)
	for {
		keys, newCursor, err := RedisClient.Scan(context.Background(), cursor, matchPattern, 500).Result()
		if err != nil {
			log.Printf("Failed to scan keys: %s\n", err.Error())
			return err
		}

		// 数据库更新当前id的喜欢列表
		for _, key := range keys {
			//先获取这个key对应的值
			videoIdList, err := RedisClient.SMembers(context.Background(), key).Result()
			if err != nil {
				log.Printf("Failed to get members: %s\n", err.Error())
				return err
			}

			//当前key对应的id
			userId, err := strconv.ParseInt(strings.Split(key, ":")[1], 10, 64)
			if err != nil {
				return err
			}
			// 点赞视频的个数
			userFavoriteCount := RedisClient.SCard(context.Background(), config.Key(config.VideoFavoriteKeyPrefix, userId)).Val()

			//创建或更新数据
			var videoIdMysql string
			for _, videoId := range videoIdList {
				videoIdMysql += videoId + ","
			}
			//不存在，就创建
			userRedisSync := domain.UserRedisSync{}
			DB.Select("user_id").Where(domain.UserRedisSync{UserId: userId}).FirstOrCreate(&userRedisSync)
			//赋值
			userRedisSync.FavoriteVideoId = videoIdMysql
			if err = DB.Model(&domain.UserRedisSync{}).
				Where("user_id = ?", userId).
				Updates(&userRedisSync).
				Error; err != nil {
				return err
			}

			//用户点赞总数跟新
			DB.Model(&domain.User{}).Where("id = ?", userId).Update("favorite_count", userFavoriteCount)
		}
		//当前这一轮redis遍历完成，下一轮
		cursor = newCursor
		if cursor == 0 {
			//迭代结束
			break
		}
	}
	return nil
}

// SyncVideoBeLikedNum 视频获赞数同步
func SyncVideoBeLikedNum() (err error) {
	//点赞视频列表同步
	matchPattern := config.VideoBeLikedNum + "*"
	cursor := uint64(0)
	for {
		keys, newCursor, err := RedisClient.Scan(context.Background(), cursor, matchPattern, 500).Result()
		if err != nil {
			log.Printf("Failed to scan keys: %s\n", err.Error())
			return err
		}
		// 数据库更新当前视频的点赞个数
		for _, key := range keys {
			//先获取这个key对应的值
			likedNum, err := RedisClient.Get(context.Background(), key).Result()
			likedNumInt64, _ := strconv.ParseInt(likedNum, 10, 64)
			if err != nil {
				log.Printf("Failed to get members: %s\n", err.Error())
				return err
			}

			//当前key对应的id
			videoId, err := strconv.ParseInt(strings.Split(key, ":")[1], 10, 64)
			if err != nil {
				return err
			}

			//更新数据
			if err = DB.Model(&domain.Video{}).Where("id", videoId).Update("FavoriteCount", likedNumInt64).Error; err != nil {
				return err
			}
		}
		//当前这一轮redis遍历完成，下一轮
		cursor = newCursor
		if cursor == 0 {
			//迭代结束
			break
		}
	}
	return nil
}

// SyncAuthorLikedNum 作者获赞数同步
func SyncAuthorLikedNum() (err error) {
	//点赞视频列表同步
	matchPattern := config.AuthorBeLikedNum + "*"
	cursor := uint64(0)
	for {
		keys, newCursor, err := RedisClient.Scan(context.Background(), cursor, matchPattern, 500).Result()
		if err != nil {
			log.Printf("Failed to scan keys: %s\n", err.Error())
			return err
		}
		// 数据库更新当前视频的点赞个数
		for _, key := range keys {
			//先获取这个key对应的值
			likedNum, err := RedisClient.Get(context.Background(), key).Result()
			likedNumInt64, _ := strconv.ParseInt(likedNum, 10, 64)
			if err != nil {
				log.Printf("Failed to get members: %s\n", err.Error())
				return err
			}

			//当前key对应的id
			userId, err := strconv.ParseInt(strings.Split(key, ":")[1], 10, 64)
			if err != nil {
				return err
			}

			//更新数据
			if err = DB.Model(&domain.User{}).Where("id", userId).Update("total_favorited", likedNumInt64).Error; err != nil {
				return err
			}
		}
		//当前这一轮redis遍历完成，下一轮
		cursor = newCursor
		if cursor == 0 {
			//迭代结束
			break
		}
	}
	return nil
}

// SyncRelationToMysql 用户关系同步
func SyncRelationToMysql() (err error) {
	//1. 同步用户关注列表到mysql中
	matchPattern := config.UserFollowHashPrefix + "*"
	cursor := uint64(0)
	for {

		err, newCursor := relationMultiplex(context.Background(), cursor, matchPattern, 300, "FollowIds")
		if err != nil {
			return err
		}
		//当前这一轮redis遍历完成,开启下一轮
		cursor = newCursor
		if cursor == 0 {
			//迭代结束
			break
		}
	}

	//2. 同步用户粉丝列表到mysql中
	matchPattern1 := config.UserFollowersHashPrefix + "*"
	cursor1 := uint64(0)
	for {
		err, nextCursor := relationMultiplex(context.Background(), cursor1, matchPattern1, 300, "FollowerIds")
		if err != nil {
			return err
		}
		//当前这一轮redis遍历完成,开启下一轮
		cursor = nextCursor
		if cursor == 0 {
			//迭代结束
			break
		}
	}
	/*err = RedisClient.Close()
	if err != nil {
		return err
	}*/
	return nil

}

// relationMultiplex 用户关注、被关注逻辑复用
func relationMultiplex(ctx context.Context, cursor uint64, matchPattern string, count int64, column string) (error, uint64) {
	keys, newCursor, err := RedisClient.Scan(ctx, cursor, matchPattern, count).Result()
	if err != nil {
		return err, 0
	}

	for _, key := range keys {
		//redis中取出每个key对应的value中的字段，不取值
		toUserIdList, err := RedisClient.HKeys(context.Background(), key).Result()
		if err != nil {
			return err, 0
		}

		//当前userId
		userId, err := strconv.ParseInt(strings.Split(key, ":")[1], 10, 64)
		if err != nil {
			return err, 0
		}

		var toUserIdMysql string
		for _, toUserId := range toUserIdList {
			toUserIdMysql += toUserId + ","
		}
		//构造数据结构
		var userRedisSync domain.UserRedisSync
		//不存在，就创建
		DB.Select("user_id").Where(domain.UserRedisSync{UserId: userId}).FirstOrCreate(&userRedisSync)
		//赋值
		if column == "FollowerIds" {
			userRedisSync.FollowerId = toUserIdMysql
		} else {
			userRedisSync.FollowId = toUserIdMysql
		}
		if err = DB.Model(&domain.UserRedisSync{}).
			Where("user_id = ?", userId).
			Updates(&userRedisSync).
			Error; err != nil {
			return err, 0
		}
	}
	return nil, newCursor
}
