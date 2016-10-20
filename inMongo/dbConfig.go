package inMongo

type row struct {
	Uid           int    `json:"uid"`
	Type          int    `json: "type"`
	Mark          int    `json: "mark"`
	Isnew         int    `json: "isnew"`
	From          int    `json: "from"`
	Channel       int    `json: "channel"`
	Channel_types int    `json: "channel_types"`
	Title         string `json: "title"`
	Content       string `json: "content"`
	Image         string `json: "image"`
	Url_type      int    `json: "url_type"`
	Url           string `json: "url"`
	Created       string `json: "created"`
	Modified      string `json: "modified"`
}
