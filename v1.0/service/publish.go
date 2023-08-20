package service

import (
	"errors"
	dao2 "github.com/goTouch/TicTok_SimpleVersion/v1.0/dao"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/domain"
	"github.com/goTouch/TicTok_SimpleVersion/v1.0/util"
	"log"
)

// InsertVideos 向数据库插入视频信息
func InsertVideos(videoName string, title string, coverName string, userId int64, coverGenerateStatus bool) error {
	video := domain.Video{
		AuthorId: userId,
		Title:    title,
		PlayUrl:  videoName,
		CoverUrl: coverName,
		//CreatTime: time.Time{}, // TODO 这个字段 为什么没赋值呢？去研究一下
	}
	// 若生成封面失败，视频的封面地址会被替换为默认封面的地址
	if !coverGenerateStatus {
		video.CoverUrl = ""
	}
	//插入数据库
	if err := dao2.DB.Create(&video).Error; err != nil {
		log.Print("向video数据库中插入数据失败！")
		log.Println(err)
		return err
	}
	//没出错
	return nil
}

// QueryAuthorPublishedVideo 查询用户发布的视频，以展示在个人列表中
func QueryAuthorPublishedVideo(authorIdInt64 int64) (videoList []domain.Video, err error) {
	url := dao2.MinioClient.EndpointURL().String() + "/" + util.VidioBucketName + "/"
	picurl := dao2.MinioClient.EndpointURL().String() + "/" + util.PictureBucketName + "/"
	err = dao2.DB.Model(&domain.Video{}).
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
