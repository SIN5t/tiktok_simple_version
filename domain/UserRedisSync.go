package domain

type UserRedisSync struct {
	UserId          int64  `json:"user_id" gorm:"primaryKey"`
	FavoriteVideoId string `json:"favorite_video_id" gorm:"type:varchar(255)"`
	FollowerId      string `json:"follower_id" gorm:"type:varchar(255)"`
	FollowId        string `json:"follow_id" gorm:"type:varchar(255)"`
}
