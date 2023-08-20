package domain

type Response struct {
	StatusCode int32  `json:"status_code"`          //结构体标记tag:`json:"status_code"` 在序列化时指定自定义的 JSON 键名
	StatusMsg  string `json:"status_msg,omitempty"` //omitempty 表示如果该字段的值为空（例如空字符串），在进行 JSON 序列化时将忽略该字段
}
