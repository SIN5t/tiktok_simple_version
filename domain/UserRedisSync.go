package domain

type UserRedisSync struct {
	UserId           int64    `json:"user_id" gorm:"primaryKey"`
	FavoriteVideoIds []string `json:"favorite_video_ids" gorm:"type:json"`
	FollowerIds      []string `json:"follower_ids" gorm:"type:json" `
	FollowIds        []string `json:"follow_ids" gorm:"type:json" `
}
