package service

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/goForward/tictok_simple_version/dao"
	"github.com/goForward/tictok_simple_version/domain"
	"github.com/goForward/tictok_simple_version/util"
	"gorm.io/gorm"
)

// Favorite
// 　业务需求：注意看官方user中，有total_favorite(当前用户获赞数量）
// 意味着，当前用户点赞的时候：视频作者获赞数量++、自己点赞数量++　TODO　favorite_count（该用户点赞数量可以用redis SCard(ctx, "myset").Result()统计)
// 在video中，有FavoriteCount、IsFavorite（由redis操作，每个用户维护自己点赞的视频列表在redis中）
func Favorite(videoIdInt64 int64, userIdInt64 int64, actionType int32) (err error) {
	///第一步一定要查userIdInt64是否合法，videoId也要查

	// TODO 先通过布隆过滤器过滤无效的用户id
	/*	if !userIdFilter.TestString(strconv.FormatInt(userIdInt64, 10)) {
		return errors.New("当前操作用户不存在")
	}*/

	//如果是点赞
	if actionType == 1 {
		//1. 在redis维护的用户点赞列表中加上该视频id
		// 先判断该用户点赞了没有
		isFavorite := dao.RedisClient.
			SIsMember(context.Background(), util.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), videoIdInt64).
			Val()
		if !isFavorite {
			//没点赞,向当前用户点赞列表中加入该视频
			dao.RedisClient.
				SAdd(context.Background(), util.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), videoIdInt64)
		}

		//2.total_favorite(当前用户获赞数量）++  使用redis做，TODO 数据库中可以不存这个字段
		//Incr 方法用于递增 Redis 中的整数值键。如果键不存在，它会将键的值初始化为 0，然后再执行增加操作
		dao.RedisClient.Incr(context.Background(), util.AuthorBeLikedNum+strconv.FormatInt(userIdInt64, 10))

		//3.video的favoriteCount
		//使用事务确保数据库更新一致性：在更新数据的过程中，数据库会自动对被修改的数据进行加锁，以防止其他并发操作同时修改同一行数据。

		//开启事务
		tx := dao.DB.Begin()
		if err := tx.Error; err != nil {
			log.Println("视频点赞：开启事务失败")
			log.Println(err)
			return errors.New("事务开启失败")
		}
		//业务逻辑
		result := dao.DB.Model(&domain.Video{}).
			Where("id = ?", videoIdInt64).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))

		if result.Error != nil {
			log.Println("数据库增加点赞数出现错误！")
			log.Println(result.Error)
		}
		if result.RowsAffected == 0 {
			log.Println("video not found")
		}

		//提交事务
		if err := tx.Commit().Error; err != nil {
			log.Println("视频点赞：事务提交失败！")
			log.Println(err)
		}

	} else if actionType == 2 { //取消点赞
		//1. 在redis维护的用户点赞列表中加上该视频id
		isFavVideo := dao.RedisClient.
			SIsMember(context.Background(), util.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), videoIdInt64).
			Val()
		if !isFavVideo { //本来就没点赞
			return errors.New("用户未曾点赞，无法取消点赞")
		}
		//点赞了现在取消
		dao.RedisClient.SRem(context.Background(), util.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), videoIdInt64)

		//2.total_favorite(当前用户获赞数量）++  使用redis做，TODO 数据库中可以不存这个字段
		//Incr 方法用于递增 Redis 中的整数值键。如果键不存在，它会将键的值初始化为 0，然后再执行增加操作
		dao.RedisClient.Decr(context.Background(), util.AuthorBeLikedNum+strconv.FormatInt(userIdInt64, 10))

		//3.video的favoriteCount
		//使用事务确保数据库更新一致性：在更新数据的过程中，数据库会自动对被修改的数据进行加锁，以防止其他并发操作同时修改同一行数据。

		//开启事务
		tx := dao.DB.Begin()
		if err := tx.Error; err != nil {
			log.Println("视频点赞：开启事务失败")
			log.Println(err)
			return errors.New("事务开启失败")
		}
		//业务逻辑
		result := dao.DB.Model(&domain.Video{}).
			Where("id = ?", videoIdInt64).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1))

		if result.Error != nil {
			log.Println("数据库增加点赞数出现错误！")
			log.Println(result.Error)
		}
		if result.RowsAffected == 0 {
			log.Println("video not found")
		}

		//提交事务
		if err := tx.Commit().Error; err != nil {
			log.Println("视频点赞：事务提交失败！")
			log.Println(err)
		}
	}
	return nil
}

func FavoriteList(userIdInt64 int64) (videoList []domain.Video, err error) {
	url := dao.MinioClient.EndpointURL().String() + "/" + util.VideoBucketName + "/"
	picurl := dao.MinioClient.EndpointURL().String() + "/" + util.PictureBucketName + "/"
	userFavoriteVideosIdStrArr, err := dao.RedisClient.
		SMembers(context.Background(), util.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10)).
		Result()
	if err != nil {
		return nil, err
	}

	for _, videoIdStr := range userFavoriteVideosIdStrArr {
		//数据库的id是int64
		videoIdInt64, err := strconv.ParseInt(videoIdStr, 10, 64)
		if err != nil {
			return nil, errors.New("字符串id解析错误")
		}
		video := domain.Video{}
		result := dao.DB.Model(&domain.Video{}).
			Where("id = ?", videoIdInt64).
			Find(&video)
		if result == nil {
			return nil, result.Error
		}
		videoList = append(videoList, video)
	}
	for i := range videoList {
		if videoList[i].CoverUrl == "" {
			videoList[i].CoverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
		} else {
			videoList[i].CoverUrl = picurl + videoList[i].CoverUrl
		}
		videoList[i].PlayUrl = url + videoList[i].PlayUrl
	}
	return videoList, nil
}
