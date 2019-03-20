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
}
