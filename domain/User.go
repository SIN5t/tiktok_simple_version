package domain

type User struct {
	IsFollow       bool   `json:"is_follow" gorm:"-"`
	Id             int64  `json:"id" gorm:"primaryKey"`
	FollowCount    int64  `json:"follow_count,omitempty" gorm:"-" `
	FollowerCount  int64  `json:"follower_count,omitempty" gorm:"-"`
	TotalFavorited int64  `json:"total_favorited,omitempty" ` //获赞次数
	FavoriteCount  int64  `json:"favorite_count,omitempty"`   //点赞个数，主页展示自己点赞视频列表的使用用到
	Salt           string `json:"-" gorm:"type:char(4)"`
	Name           string `json:"name" gorm:"type:varchar(32); index:idx_usr_name"`
	Pwd            string `json:"-" gorm:"type:char(60)"`
	Avatar         string `json:"avatar" gorm:"type:char(60)"` //头像
	WorkCount      int64  `json:"work_count,omitempty"`        //作品数量，主页展示作品个数用到
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type UserFollowListResponse struct {
	Response
	UserFollowList []User `json:"user_list"` //标签（Tag），用于指定结构体字段在序列化为JSON格式时的命名规则,之前漏写，导致前端无法识别！
}
