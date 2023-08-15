package domain

import "time"

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}
type Message struct {
	Id         int64     `json:"id"  gorm:"primaryKey"`
	ToUserId   int64     `json:"to_user_id" `
	FromUserId int64     `json:"from_user_id" `
	Content    string    `json:"content" gorm:"type:varchar(200)"`
	CreateTime time.Time `json:"create_time" gorm:"type:datetime(0);autoCreateTime;index:;sort:desc"`
}
