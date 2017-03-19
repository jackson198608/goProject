package post

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/squirrel"
	"github.com/jackson198608/structable"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

// For convenience, we declare the table name as a constant.
var baseTable string = "pre_forum_post"
var tableBaseNum int = 100

// This is our struct. Notice that we make this a structable.Recorder.
type Post struct {
	structable.Recorder
	builder squirrel.StatementBuilderType

	//Pid         int64  `stbl:"pid,PRIMARY_KEY,SERIAL"`
	Pid         int64  `stbl:"pid"`
	Fid         int64  `stbl:"fid"`
	Tid         int64  `stbl:"tid"`
	First       int    `stbl:"first"`
	Author      string `stbl:"author"`
	Authorid    int64  `stbl:"authorid"`
	Subject     string `stbl:"subject"`
	Dateline    int64  `stbl:"dateline"`
	Message     string `stbl:"message"`
	Useip       string `stbl:"useip"`
	Invisible   int    `stbl:"invisible"`
	Anonymous   int    `stbl:"anonymous"`
	Usesig      int    `stbl:"usesig"`
	Htmlon      int    `stbl:"htmlon"`
	Bbcodeoff   int    `stbl:"bbcodeoff"`
	Smileyoff   int    `stbl:"smileyoff"`
	Parseurloff int    `stbl:"parseurloff"`
	Attachment  int    `stbl:"attachment"`
	Rate        int    `stbl:"rate"`
	Ratetimes   int    `stbl:"ratetimes"`
	Status      int    `stbl:"status"`
	Tags        string `stbl:"tags"`
	Comment     int    `stbl:"comment"`
	Replycredit int    `stbl:"replycredit"`
	//Position    int64  `stbl:"position"`
	isSplit   bool
	logLevel  int
	postTable string
}

// NewUser creates a new Structable wrapper for a user.
//
// Of particular importance, watch how we intialize the Recorder.
func NewPost(logLevel int, db squirrel.DBProxyBeginner, dbFlavor string, pid int64, tid int64, isSplit bool) *Post {
	u := new(Post)
	//logger.SetRollingDaily("/tmp", "1.log")
	logger.SetLevel(logger.LEVEL(logLevel))

	u.isSplit = isSplit
	if (pid > 0) && (tid > 0) {
		u.Pid = pid
		u.Tid = tid
	}

	u.postTable = baseTable

	if isSplit && (tid > 0) {
		u.postTable = u.getTableSplitName()
	}
	logger.Info(u.postTable)
	u.Recorder = structable.New(db, dbFlavor).Bind(u.postTable, u)

	if (pid > 0) && (tid > 0) {
		u.LoadByPid()
	}

	u.logLevel = logLevel

	return u
}

func (p *Post) PidExists() bool {
	isExist, err := p.ExistsWhere("pid = ?", p.Pid)
	if err != nil {
		logger.Error("find exists error", p.Tid, p.Pid, p.TableName(), err)
	}
	return isExist
}
func (p *Post) hasChanged() bool {
	if p.Pid <= 0 || p.Tid <= 0 {
		logger.Error("have no pid or tid can not continute")
		return false
	}
	if !p.isSplit {
		p.postTable = p.getTableSplitName()
		p.Recorder.ChangeBindTableName(p.postTable)
		p.isSplit = true
	}
	isExist := p.PidExists()

	p.backToMain()
	return isExist

}

func (p *Post) backToMain() bool {
	p.postTable = baseTable
	p.isSplit = false
	p.Recorder.ChangeBindTableName(p.postTable)
	return true
}

func (p *Post) MoveToSplit() bool {
	if p.hasChanged() {
		logger.Info("has changed", p.Pid, p.Tid)
		return true
	} else {
		p.postTable = p.getTableSplitName()
		p.Recorder.ChangeBindTableName(p.postTable)
		defer p.backToMain()
		p.isSplit = true
		err := p.Insert()
		if err != nil {
			logger.Error("insert error", p.Pid, p.Tid, p.TableName(), err)
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
	PostTableSplit := baseTable + "_" + tableNumStr
	return PostTableSplit
}

// LoadByName is a custom loader.
//
// The Load() method on a Recorder loads by ID. This allows us to load by
// a different field -- Name.
func (p *Post) LoadByPid() error {
	return p.Recorder.LoadWhere("pid = ? limit 0,1", p.Pid)
}
