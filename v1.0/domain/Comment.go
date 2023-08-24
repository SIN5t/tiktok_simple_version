package domain

type Comment struct {
	Id         int64  `json:"id" gorm:"primaryKey"`
	UserId     int64  `json:"-"`
	VideoId    int64  `json:"video_id,omitempty"`
	CreateDate string `json:"create_date" gorm:"type:varchar(10);index"`
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
