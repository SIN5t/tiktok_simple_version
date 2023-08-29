package domain

type UserRedisSync struct {
	Id               int64    `json:"id" gorm:"primaryKey"`
	FavoriteVideoIds []string `json:"favorite_video_ids" gorm:"type:json"`
	FollowerIds      []string `json:"follower_ids" gorm:"type:json" `
	FollowIds        []string `json:"follow_ids" gorm:"type:json" `
}
