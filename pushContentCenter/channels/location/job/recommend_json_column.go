package job

//json column
type RecommendJsonColumn struct {
	Infoid        int
	Pid           int
	Uid           int
	Ruid          int
	Type          int
	Tag           int    //帖子的热门话题ID
	Tags          string //标签
	QstType       int
	AdType        int
	AdUrl         string
	Title         string
	Description   string
	Images        string
	Imagenums     int
	Created       int
	Action        int
	Channel       int
	Rauth         string //认证信息
	RecommendType string
}
