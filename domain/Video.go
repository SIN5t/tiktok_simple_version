package domain

import "time"

type Video struct {
	IsFavorite    bool      `json:"is_favorite" gorm:"-"`
	Id            int64     `json:"id" gorm:"primaryKey"`
	AuthorId      int64     `json:"author_id"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	Title         string    `json:"title" gorm:"type:varchar(100)"`
	PlayUrl       string    `json:"play_url" gorm:"type:varchar(100)"`
	CoverUrl      string    `json:"cover_url" gorm:"type:varchar(100)"`
	CreatTime     time.Time `json:"-" gorm:"type:datetime(0);autoCreateTime;index:;sort:desc"` //该字段加了索引

	Author User `json:"author"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
