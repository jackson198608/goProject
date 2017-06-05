package task

import (
	"database/sql"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/eventLog/Pushdata"
	mgo "gopkg.in/mgo.v2"
	"strconv"
	"strings"
)

type Task struct {
	loggerLevel int
	id          int
	oid         int
	uid         int
	fuid        int
	fans        string //粉丝数
	ecount      string //动态数
	status      string
	db          *sql.DB
	session     *mgo.Session
	// event       *EventLogNew
	count      int
	fansnum    int
	loopNum    int
	fansLimit  int
	eventLimit int
	pushLimit  int
	dateLimit  string
}
type Follow struct {
	follow_id int
}

func NewTask(loggerLevel int, redisStr string, db *sql.DB, session *mgo.Session) *Task {
	if loggerLevel < 0 {
		loggerLevel = 0
	}
	logger.SetLevel(logger.LEVEL(loggerLevel))

	redisArr := strings.Split(redisStr, "|")
	var fuid int
	var uid int
	var id int
	var oid int
	var status string
	var fans string
	var ecount string
	if len(redisArr) == 2 {
		if redisArr[1] == "3" {
			uids := strings.Split(redisArr[0], "&")
			fuid, _ = strconv.Atoi(uids[0])
			uid, _ = strconv.Atoi(uids[1])
		} else {
			id, _ = strconv.Atoi(redisArr[0])
			status = redisArr[1] //要执行的操作:0:删除,-1隐藏,1显示,2动态推送给粉丝
		}
	}

	if len(redisArr) == 1 {
		oid, _ = strconv.Atoi(redisStr)
	}

	if len(redisArr) == 3 {
		uid, _ = strconv.Atoi(redisArr[0]) //用户uid
		fans = redisArr[1]                 //粉丝数
		ecount = redisArr[2]               //动态数
	}
	t := new(Task)
	t.oid = oid
	t.id = id
	t.fuid = fuid
	t.uid = uid
	t.status = status
	t.session = session
	t.db = db
	t.fans = fans
	t.ecount = ecount
	return t

}

func (t *Task) Do() {
	m := Pushdata.NewEventLogNew(t.loggerLevel, t.oid, t.id, t.db, t.session)
	if m != nil {
		if t.oid > 0 {
			logger.Info("export event to mongo")
			m.SaveMongoEventLog(t.oid)
		}
		if t.id > 0 {
			logger.Info("update mongo event status")
			m.UpdateMongoEventLogStatus(t.id, t.status)
		}
		if t.fuid > 0 && t.uid > 0 {
			logger.Info("remove fans event")
			m.RemoveFansEventLog(t.fuid, t.uid)
		}
	}
}

func (t *Task) Dopush(dateLimit string, loopNum int, fansLimit string, eventLimit string, pushLimit string) {
	m := Pushdata.NewEventLogNew(t.loggerLevel, t.oid, t.id, t.db, t.session)
	m.PushEventToFansTask(t.fans, t.uid, t.ecount, loopNum, fansLimit, eventLimit, pushLimit, dateLimit)

}
