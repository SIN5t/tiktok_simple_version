package service

import (
	"context"
	"errors"
	"github.com/goForward/tictok_simple_version/config"
	"github.com/goForward/tictok_simple_version/dao"
	"github.com/goForward/tictok_simple_version/domain"
	"strconv"
	"time"
)

// Favorite 点赞、取消接口
func Favorite(videoIdInt64 int64, userIdInt64 int64, actionType int32) (err error) {
	// TODO 先通过布隆过滤器过滤无效的用户id
	/*	if !userIdFilter.TestString(strconv.FormatInt(userIdInt64, 10)) {
		return errors.New("当前操作用户不存在")
	}*/

	// 找到作者id
	var authorId int64
	dao.DB.Model(&domain.Video{}).Where("Id = ?", videoIdInt64).Select("author_id").Find(&authorId)

	//点赞
	if actionType == 1 {
		//1. 在redis维护的用户点赞列表中加上该视频id
		// 判断是否点赞
		isFavorite := dao.RedisClient.
			SIsMember(context.Background(), config.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), videoIdInt64).
			Val()
		if !isFavorite {
			//没点赞,向当前用户点赞列表中加入该视频
			dao.RedisClient.
				SAdd(context.Background(), config.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), videoIdInt64)
			dao.RedisClient.Expire(context.Background(), config.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), 6*30*24*time.Hour)
		}

		//2.total_favorite(点赞视频对应的作者获赞数量增加），dao处有定时同步到mysql的逻辑

		//如果键不存在，它会将键的值初始化为 0，然后再执行增加操作
		dao.RedisClient.Incr(context.Background(), config.AuthorBeLikedNum+strconv.FormatInt(authorId, 10))
		dao.RedisClient.Expire(context.Background(), config.AuthorBeLikedNum+strconv.FormatInt(authorId, 10), 6*30*24*time.Hour)
		//3.video的favoriteCount，dao中有逻辑同步到mysql表中
		dao.RedisClient.Incr(context.Background(), config.VideoBeLikedNum+strconv.FormatInt(videoIdInt64, 10))
		dao.RedisClient.Expire(context.Background(), config.VideoBeLikedNum+strconv.FormatInt(videoIdInt64, 10), 6*30*24*time.Hour)
		/*//开启事务
		tx := dao.DB.Begin()
		if err := tx.Error; err != nil {
			return err
		}
		//视频获赞+1
		if result := tx.Model(&domain.Video{}).
			Where("id = ?", videoIdInt64).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).
			Error; result != nil {
			tx.Rollback()
			return err
		}
		//视频作者获赞数+1
		//	通过视频id找到视频对应的作者
		tx.Model(&domain.Video{}).Where("Id = ?", videoIdInt64).Select("author_id").Find(&authorId)
		if err = tx.Model(&domain.User{}).
			Where("AuthorId = ?", authorId).
			UpdateColumn("TotalFavorited", gorm.Expr("total_favorited + ?", 1)).
			Error; err != nil {
			tx.Rollback()
			return err
		}

		//提交事务
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}*/

	} else if actionType == 2 { //取消点赞
		//1. 在redis维护的用户点赞列表中加上该视频id
		isFavVideo := dao.RedisClient.
			SIsMember(context.Background(), config.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), videoIdInt64).
			Val()
		if !isFavVideo { //本来就没点赞
			return errors.New("用户未曾点赞，无法取消点赞")
		}
		//取消点赞
		dao.RedisClient.SRem(context.Background(), config.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10), videoIdInt64)

		//2.total_favorite(当前视频作者获赞数量）
		dao.RedisClient.Decr(context.Background(), config.AuthorBeLikedNum+strconv.FormatInt(authorId, 10))

		//3.video的favoriteCount
		dao.RedisClient.Decr(context.Background(), config.VideoBeLikedNum+strconv.FormatInt(videoIdInt64, 10))

		/*//开启事务
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
		}*/
	}
	return nil
}

func FavoriteList(userIdInt64 int64) (videoList []domain.Video, err error) {
	url := dao.MinioClient.EndpointURL().String() + "/" + config.VideoBucketName + "/"
	picurl := dao.MinioClient.EndpointURL().String() + "/" + config.PictureBucketName + "/"
	userFavoriteVideosIdStrArr, err := dao.RedisClient.
		SMembers(context.Background(), config.VideoFavoriteKeyPrefix+strconv.FormatInt(userIdInt64, 10)).
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
			Select("id,author_id,title,cover_url,favorite_count").
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
