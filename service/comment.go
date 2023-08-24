package service

import (
	"errors"
	"time"

	"github.com/goForward/tictok_simple_version/dao"
	"github.com/goForward/tictok_simple_version/domain"
	"gorm.io/gorm"
)

func Comment(commentId int64) (comment domain.Comment, err error) {
	if commentId <= 0 {
		return comment, errors.New("不合法的评论id")
	}
	err = dao.DB.Model(&domain.Comment{}).Where("id = ?", commentId).Find(&comment).Error
	if err != nil {
		return domain.Comment{}, errors.New("评论不存在")
	}
	return comment, nil
}

func AddComment(videoId, userId int64, content string) (comment domain.Comment, err error) {
	// 判断用户是否存在
	user, err := User(userId)
	if err != nil {
		return comment, err
	}

	comment = domain.Comment{
		UserId:     userId,
		VideoId:    videoId,
		CreateDate: time.Now().Format("01-02"),
		Content:    content,
		User:       user,
	}

	// 创建评论事务
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		video := domain.Video{}

		//判断视频是否存在
		err = tx.Model(&video).Where("id = ?", videoId).Find(&video).Error
		if err != nil {
			return err
		}

		// 增加评论数量
		commentCount := video.CommentCount + 1
		err = tx.Model(&video).Where("id = ?", videoId).Update("comment_count", &commentCount).Error
		if err != nil {
			return err
		}

		err = tx.Model(&domain.Comment{}).Create(&comment).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Comment{}, err
	}

	return comment, nil
}

func DeleteComment(commentId int64) (err error) {
	// 评论是否存在
	comment, err := Comment(commentId)
	if err != nil {
		return err
	}

	// 删除评论事务
	videoId := comment.VideoId
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		video := domain.Video{}

		// 判断视频是否存在
		err = tx.Model(&video).Where("id = ?", videoId).Find(&video).Error
		if err != nil {
			return err
		}

		// 减少评论数量
		commentCount := video.CommentCount - 1
		err = tx.Model(&video).Where("id = ?", videoId).Update("comment_count", commentCount).Error
		if err != nil {
			return err
		}

		// 删除评论
		err = tx.Model(&domain.Comment{}).Delete(&comment).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func CommentList(videoId int64) (commentList []domain.Comment, err error) {
	err = dao.DB.Model(&domain.Comment{}).Where("video_id = ?", videoId).Find(&commentList).Error
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(commentList); i++ {
		comment := &commentList[i]
		user, err := User(comment.UserId)
		if err != nil {
			return nil, err
		}
		comment.User = user
	}

	return commentList, nil
}
