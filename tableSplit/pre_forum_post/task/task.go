package task

import (
	"database/sql"
	"errors"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackson198608/goProject/tableSplit/pre_forum_post/post"
	"github.com/jackson198608/squirrel"
	"strconv"
)

type Task struct {
	Tid      int64
	logLevel int
	dbAuth   string
	dbDsn    string
	dbName   string
	con      *sql.DB
	dbCache  squirrel.DBProxyBeginner
	pids     []int64
}

func NewTask(logLevel int, tidStr string, args []string) *Task {

	logger.SetLevel(logger.LEVEL(logLevel))

	//check the string
	if len(args) != 3 {
		logger.Error("there is not enough args to start")
		return nil
	}

	//set value
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		logger.Error("tid error")
		return nil
	}

	t := new(Task)
	t.Tid = int64(tid)
	t.dbAuth = args[0]
	t.dbDsn = args[1]
	t.dbName = args[2]
	t.logLevel = logLevel

	//make db comon value and check error
	err = t.getDbCache()
	if err != nil {
		logger.Error("get pid list error", err, t.Tid)
		return nil
	}

	err = t.getPids()
	if err != nil {
		logger.Error("get pid list error", err, t.Tid)
		t.con.Close()
		return nil
	}

	//if have no pid,no task
	if len(t.pids) == 0 {
		logger.Error("this tid have no pid ,so pass", t.Tid)
		t.con.Close()
		return nil
	}

	return t
}

func (t *Task) getDbCache() error {
	con, err := sql.Open("mysql", t.dbAuth+"@tcp("+t.dbDsn+")/"+t.dbName+"?charset=utf8")
	if err != nil {
		logger.Error("connect err", t.dbDsn, t.dbAuth, t.dbName)
		return errors.New("connect db error")

	}
	// Third, we wrap in a prepared statement cache for better performance.
	cache := squirrel.NewStmtCacheProxy(con)
	t.dbCache = cache
	t.con = con
	return nil
}

func (t *Task) Do() error {
	if t.pids == nil {
		return errors.New("pids is nil")
	}
	pidlen := len(t.pids)
	for i := 0; i < pidlen; i++ {
		t.handlePid(t.pids[i])
	}
	return nil
}

func (t *Task) handlePid(pid int64) error {
	if pid <= 0 {
		logger.Error("pid <0")
		return errors.New("pid<0")
	}
	post := post.NewPost(t.logLevel, t.dbCache, "mysql", pid, t.Tid, false)
	exist := post.PidExists()
	if exist {
		logger.Info("doing for this pid", post.Message, pid)
		result := post.MoveToSplit()
		if !result {
			logger.Error("change pid error", pid)
			return errors.New("change pid error")
		}
	} else {
		logger.Info("not exsit", pid)
		return nil
	}

	return nil

}

func (t *Task) getPids() error {
	//connect db
	db, err := sql.Open("mysql", t.dbAuth+"@tcp("+t.dbDsn+")/"+t.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("connect db error", t.dbDsn, t.dbAuth, t.dbName)

		//fmt.Println("[error] connect db err")
		return errors.New("connect db error")
	}
	defer db.Close()

	//do query part
	rows, err := db.Query("select pid from pre_forum_post where tid=" + strconv.Itoa(int(t.Tid)))
	if err != nil {
		logger.Error("get pids by tid error", t.Tid, err)
		return errors.New("get pids error")
	}

	//get what we need

	pids := make([]int64, 0, 100)
	for rows.Next() {
		var pid int64
		if err := rows.Scan(&pid); err != nil {
			logger.Error("get pid error after row.next", err)
			return errors.New("get pid error after row.next")
		}
		logger.Info("find pid:", pid, "in tid:", t.Tid)
		pids = append(pids, pid)
	}
	if err := rows.Err(); err != nil {
		logger.Error("rows error agrain ", err)
		return errors.New("rows error agrain")
	}
	t.pids = pids
	return nil
}
