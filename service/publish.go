package service

import (
	"errors"
	"github.com/goForward/tictok_simple_version/config"
	"gorm.io/gorm"
	"log"

	"github.com/goForward/tictok_simple_version/dao"
	"github.com/goForward/tictok_simple_version/domain"
)

// InsertVideos 向数据库插入视频信息
func InsertVideos(videoName string, title string, coverName string, userId int64, coverGenerateStatus bool) error {
	snowFakeId := dao.VideoNode.Generate().Int64()
	video := domain.Video{
		Id:       snowFakeId,
		AuthorId: userId,
		Title:    title,
		PlayUrl:  videoName,
		CoverUrl: coverName,
		//CreatTime: time.Time{},
	}
	// 若生成封面失败，视频的封面地址会被替换为默认封面的地址
	if !coverGenerateStatus {
		video.CoverUrl = ""
	}
	tx := dao.DB.Begin()
	//插入数据库
	if err := tx.Create(&video).Error; err != nil {
		log.Print("向video数据库中插入数据失败！")
		tx.Rollback()
		return err
	}
	//用户信息那边发布视频数量+1
	if err := tx.Model(&domain.User{}).
		Where("id = ?", userId).
		UpdateColumn("work_count", gorm.Expr("work_count + ?", 1)).
		Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// QueryAuthorPublishedVideo 查询用户发布的视频，以展示在个人列表中
func QueryAuthorPublishedVideo(authorIdInt64 int64) (videoList []domain.Video, err error) {
	url := dao.MinioClient.EndpointURL().String() + "/" + config.VideoBucketName + "/"
	picurl := dao.MinioClient.EndpointURL().String() + "/" + config.PictureBucketName + "/"
	err = dao.DB.Model(&domain.Video{}).
		Select("id,favorite_count,cover_url,play_url").
		Where("author_id = ?", authorIdInt64).
		Order("creat_time desc"). //该字段加了索引
		Find(&videoList).Error
	if err != nil {
		err = errors.New("未查询到视频")
		return
	}
	for i := range videoList {
		if videoList[i].CoverUrl == "" {
			videoList[i].CoverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
		} else {
			videoList[i].CoverUrl = picurl + videoList[i].CoverUrl
		}

		videoList[i].PlayUrl = url + videoList[i].PlayUrl
	}
	return
}
