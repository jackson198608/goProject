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
	fans        []*Follow
	status      string
	db          *sql.DB
	session     *mgo.Session
	// event       *EventLogNew
	count   int
	fansnum int
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
	if len(redisArr) == 2 {
		if redisArr[1] == "3" {
			uids := strings.Split(redisArr[1], "&")
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

	t := new(Task)
	t.oid = oid
	t.id = id
	t.fuid = fuid
	t.uid = uid
	t.status = status
	t.session = session
	t.db = db
	// t.event = LoadById(id, db)
	return t

}

func (t *Task) Do() {
	m := Pushdata.NewEventLogNew(t.loggerLevel, t.oid, t.id, t.db, t.session)
	if m != nil {
		if t.oid > 0 {
			m.SaveMongoEventLog(t.oid)
		}
		if t.id > 0 {
			m.UpdateMongoEventLogStatus(t.id, t.status)
		}
		if t.fuid > 0 && t.uid > 0 {
			m.RemoveFansEventLog(t.fuid, t.uid)
		}
	}
}

// queueName := t.queueName //+ "_" + tableNumStr
//     logger.Info("pop ", queueName)
//     db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
//     if err != nil {
//         logger.Error("[error] connect db err")
//     }
//     defer db.Close()
//     session, err := mgo.Dial(c.mongoConn)
//     if err != nil {
//         return
//     }
//     defer session.Close()
//     for {
//         //doing until got nothing
//         redisStr := (*t.client).LPop(queueName).Val()
//         if redisStr == "" {
//             logger.Info("got nothing", queueName)
//             x <- 1
//             return
//         }
//         redisArr := strings.Split(redisStr, "|")
//         if len(redisArr) == 2 {
//             if redisArr[1] == "3" {
//                 uids := strings.Split(redisArr[0], "&")
//                 RemoveFansEventLog(uids[0], uids[1], session) //fuid,uid
//             } else {
//                 u := LoadMongoById(redisArr[0], session)
//                 fans := GetFansData(u.Uid, db)
//                 status := redisArr[1] //要执行的操作:0:删除,-1隐藏,1显示,2动态推送给粉丝
//                 UpdateMongoEventLogStatus(u, fans, status, session)
//             }
//         }
//         //doing job
//         if len(redisArr) == 1 {
//             id, _ := strconv.Atoi(redisStr)
//             u := LoadById(id, db)
//             // fans := GetFansData(u.uid, db)
//             var fans []*Follow
//             SaveMongoEventLog(u, fans, session)
//         }
//     }
