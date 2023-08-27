package dao

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/goForward/tictok_simple_version/domain"
	"github.com/goForward/tictok_simple_version/util"
	"log"
	"strconv"
	"strings"
	"time"
)

//数据库与redis一致性同步处理

func ScheduleSyncFavoriteToMysql() error {
	scheduler := gocron.NewScheduler(time.Local) ///定时任务
	_, err := scheduler.Every(10).
		Tag("favoriteRedis").
		Seconds().
		Do(SyncFavoriteToMysql)
	if err != nil {
		return err
	}
	scheduler.StartAsync() //异步执行
	return nil
}

func SyncFavoriteToMysql() (err error) {
	//1. 视频角度： 被点赞个数，先于是否使用redis

	//2.  用户角度： 该用户点赞的视频id、该用户点赞的总数（前端没做）
	// 2.1 遍历所有相关的redis key，对每一个redis中存储的用户 的点赞视频进行同步数据库
	matchPattern := util.VideoFavoriteKeyPrefix + "*"
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
			//TODO go无法存储[]string格式，求出的时候需要反序列化
			if len(videoIdList) == 0 {
				return nil
			}
			if err != nil {
				log.Printf("Failed to get members: %s\n", err.Error())
				return err
			}
			//这个key对应的id
			userId, err := strconv.ParseInt(strings.Split(key, ":")[1], 10, 64)
			if err != nil {
				return err
			}
			//更新数据库
			if err = DB.
				Model(&domain.User{}).
				Where("Id = ?", userId).
				UpdateColumn("FavoriteVideoIds", videoIdList).
				Error; err != nil {
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
	/*err = RedisClient.Close()
	if err != nil {
		return err
	}*/
	return nil
}

func ScheduleSyncRelationToMysql() error {
	scheduler := gocron.NewScheduler(time.Local)
	_, err := scheduler.Every(10).
		Tag("relationRedis").
		Second().
		Do(SyncRelationToMysql)
	if err != nil {
		return err
	}
	scheduler.StartAsync()
	return nil
}

func SyncRelationToMysql() (err error) {
	//1. 同步用户关注列表到mysql中
	matchPattern := util.UserFollowHashPrefix + "*"
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
	matchPattern1 := util.UserFollowersHashPrefix + "*"
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

func relationMultiplex(ctx context.Context, cursor uint64, matchPattern string, count int64, column string) (error, uint64) {
	keys, newCursor, err := RedisClient.Scan(ctx, cursor, matchPattern, count).Result()
	if err != nil {
		return err, 0
	}

	for _, key := range keys {
		//redis中取出每个key对应的value中的字段，不取值
		toUserIdList, err := RedisClient.HKeys(context.Background(), key).Result()
		if len(toUserIdList) == 0 {
			return nil, 0
		}
		if err != nil {
			return err, 0
		}
		//当前userId
		userId, err := strconv.ParseInt(strings.Split(key, ":")[1], 10, 64)
		if err != nil {
			return err, 0
		}
		//更新到数据库
		if err = DB.
			Model(&domain.User{}).
			Where("Id = ?", userId).
			//Update("FollowIds", toUserIdList).
			Update(column, toUserIdList).
			Error; err != nil {
			return err, 0
		}

	}
	return nil, newCursor

}
