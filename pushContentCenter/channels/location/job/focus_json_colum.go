package job

//json column
type FocusJsonColumn struct {
	TypeId    int
	Uid       int
	Created   string
	Infoid    int
	Status    int
	Tid       int
	Bid       int
	Fid       int
	Content   string
	Title     string
	Imagenums int
	ImageInfo string
	VideoUrl  string
	Duration  int
	IsVideo   int
	Forum     string
	Tag       int
	Qsttype   int
	Source    int
	Action    int
	PetId     int //星球卡片ID
	PetType   int //宠物类型 1猫 2狗
	AdoptId int
	PetName string
	PetAge string
	PetBreed int
	PetGender int
	PetSpecies string
	Province string
	City string
	County string
	Reason string
	Image string
	PetImmunity int
	PetExpelling int
	PetSterilization int
	PetStatus int
	AdoptStatus int
	PetIntroduction string
	UserIdentity int
	AdoptTag interface{}
	PetAgenum int
	Channel int  //1 水军数据
	RegisterTime int  //注册时间
	ReplyNum     int  //回复数
	Inhome   int
	IsDigest int //是否精华
	Pid   int  //楼层id
	TagInfo  string //标签信息
	ImageWidth  int  //图片宽度
	ImageHeight int  //图片高度
	ThreadStatus int  //pre_forum_thread表中的status
}
