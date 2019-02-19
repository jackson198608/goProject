package job

type MsgPushRecordJsonColumn struct {
	Uid int
	Type int   //消息类型
	Mark int   //MessagePush主键
	Isnew int
	From  int  //来源，手机类型  0ios  1android
	Channel  int
	ChannelTypes  int
	Title  string
	Content  string
	Image  string
	UrlType  int
	Url  string
	Created  string
	Modified string
}