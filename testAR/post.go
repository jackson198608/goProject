package post

import (
	"github.com/Masterminds/squirrel"
	"github.com/Masterminds/structable"
	//_ "github.com/lib/pq"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

// For convenience, we declare the table name as a constant.
var baseTable string = "pre_forum_post"
var PostTable string = "pre_forum_post"
var tableBaseNum int = 10

// This is our struct. Notice that we make this a structable.Recorder.
type Post struct {
	structable.Recorder
	builder squirrel.StatementBuilderType

	//Pid         int64  `stbl:"pid,PRIMARY_KEY,SERIAL"`
	Pid         int64  `stbl:"pid"`
	Fid         int64  `stbl:"fid"`
	Tid         int64  `stbl:"tid"`
	First       bool   `stbl:"first"`
	Author      string `stbl:"author"`
	Authorid    int64  `stbl:"authorid"`
	Subject     string `stbl:"subject"`
	Dateline    int64  `stbl:"dateline"`
	Message     string `stbl:"message"`
	Useip       string `stbl:"useip"`
	Invisible   bool   `stbl:"invisible"`
	Anonymous   bool   `stbl:"anonymous"`
	Usesig      bool   `stbl:"usesig"`
	Htmlon      bool   `stbl:"htmlon"`
	Bbcodeoff   bool   `stbl:"bbcodeoff"`
	Smileyoff   bool   `stbl:"smileyoff"`
	Parseurloff bool   `stbl:"parseurloff"`
	Attachment  bool   `stbl:"attachment"`
	Rate        int    `stbl:"rate"`
	Ratetimes   int    `stbl:"ratetimes"`
	Status      int    `stbl:"status"`
	Tags        string `stbl:"tags"`
	Comment     bool   `stbl:"comment"`
	Replycredit int    `stbl:"replycredit"`
	Position    int64  `stbl:"position"`
	isSplit     bool
}

// NewUser creates a new Structable wrapper for a user.
//
// Of particular importance, watch how we intialize the Recorder.
func NewPost(db squirrel.DBProxyBeginner, dbFlavor string, pid int64, tid int64, isSplit bool) *Post {
	u := new(Post)
	u.isSplit = isSplit
	if (pid > 0) && (tid > 0) {
		u.Pid = pid
		u.Tid = tid
	}

	if isSplit && (tid > 0) {
		PostTable = u.getTableSplitName()
	}
	u.Recorder = structable.New(db, dbFlavor).Bind(PostTable, u)

	if (pid > 0) && (tid > 0) {
		u.LoadByPid()
	}

	return u
}

func (p *Post) PidExists() bool {
	isExist, err := p.ExistsWhere("pid = ?", p.Pid)
	if err != nil {
		fmt.Println("[error] find exists error", err)
	}
	return isExist
}
func (p *Post) hasChanged() bool {
	if p.Pid <= 0 || p.Tid <= 0 {
		fmt.Println("[error] have no pid or tid can not continute")
		return false
	}
	if !p.isSplit {
		PostTable = p.getTableSplitName()
		p.Recorder.ChangeBindTableName(PostTable)
		p.isSplit = true
	}
	isExist := p.PidExists()

	p.backToMain()
	return isExist

}

func (p *Post) backToMain() bool {
	PostTable = baseTable
	p.isSplit = false
	p.Recorder.ChangeBindTableName(PostTable)
	return true
}

func (p *Post) MoveToSplit() bool {
	if p.hasChanged() {
		fmt.Println("[info] has changed", p.Pid, p.Tid)
		return true
	} else {

		PostTable = p.getTableSplitName()
		fmt.Println("[info] former pid", p.Pid)
		p.Recorder.ChangeBindTableName(PostTable)
		fmt.Println("[info] after pid", p.Pid)
		defer p.backToMain()
		p.isSplit = true
		err := p.Insert()
		if err != nil {
			fmt.Println("[error] insert error", p.Pid, p.Tid, err)
			return false
		}
		return true

	}
}

func (p *Post) getTableSplitName() string {
	tableNum := p.Tid % int64(tableBaseNum)
	if tableNum == 0 {
		tableNum = int64(tableBaseNum)
	}
	tableNumStr := strconv.Itoa(int(tableNum))
	PostTableSplit := PostTable + "_" + tableNumStr
	return PostTableSplit
}

// LoadByName is a custom loader.
//
// The Load() method on a Recorder loads by ID. This allows us to load by
// a different field -- Name.
func (p *Post) LoadByPid() error {
	return p.Recorder.LoadWhere("pid = ? limit 0,1", p.Pid)
}
