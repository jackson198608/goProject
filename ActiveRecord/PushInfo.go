package ActiveRecord

type PushInfo struct {
	Type    int    `json:"type"`
	Mark    string `json:"mark"`
	Uid     int    `json:"uid"`
	Title   string `json:"title"`
	Content string `json:"content"`
}