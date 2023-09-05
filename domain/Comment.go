package domain

type Comment struct {
	Id         int64  `json:"id" gorm:"primaryKey"`
	UserId     int64  `json:"-" gorm:"index:idx_comm_usrId"`
	VideoId    int64  `json:"video_id,omitempty" gorm:"index:idx_comm_videoId"`
	CreateDate string `json:"create_date" gorm:"type:varchar(10);index:idx_comm_date"`
	Content    string `json:"content" gorm:"type:text"`

	User User `json:"user" gorm:"-"`
}

type CommentResponse struct {
	Response
	Comment Comment `json:"comment"`
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"`
}
