package dao

import (
	"context"
	"gorm.io/plugin/dbresolver"
	"time"
)

// Video TODO 查一下是都不用自动生成的gen下的video的struct吗？这样到时候还要类型转换一下
type Video struct {
	ID            int64     `gorm:"primarykey"`
	CreatedAt     time.Time `gorm:"not null;index:idx_create" json:"created_at,omitempty"`
	AuthorID      int64     `gorm:"index:idx_authorid;not null" json:"author_id,omitempty"`
	PlayUrl       string    `gorm:"type:varchar(255);not null" json:"play_url,omitempty"`
	CoverUrl      string    `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount uint      `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	CommentCount  uint      `gorm:"default:0;not null" json:"comment_count,omitempty"`
	Title         string    `gorm:"type:varchar(50);not null" json:"title,omitempty"`
}

func GetVideosByLastTime(ctx context.Context, lastTime int64, limit int) (videoList []*Video, err error) {

	//Clauses为数据库查询添加额外的选项或条件,WithContext可以做：超时和取消、事务、日志等
	conn := GetDB().Clauses(dbresolver.Read).WithContext(ctx)
	//如果传过来的lastTime是0，或者没传，都应该要给出小于当前时间的视频
	if lastTime == 0 {
		lastTime = time.Now().UnixMilli() //用现在的时间来替换
	}

	lastTimeUnixMilli := time.UnixMilli(lastTime)
	err = conn.Limit(limit).Order("create_at desc").Where("create_at < ?", lastTimeUnixMilli).Find(&videoList).Error
	if err != nil {
		return nil, err
	}

	return videoList, nil
}
