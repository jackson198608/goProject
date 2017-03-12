package task

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackson198608/goProject/tableSplit/pre_forum_post/post"
	"github.com/jackson198608/squirrel"
	"strconv"
)

type Task struct {
	Tid     int64
	dbAuth  string
	dbDsn   string
	dbName  string
	con     *sql.DB
	dbCache squirrel.DBProxyBeginner
	pids    []int64
}

func NewTask(tidStr string, dbAuth string, dbDsn string, dbName string) *Task {
	//set value
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		fmt.Println("[error]tid error")
		return nil
	}
	t := new(Task)
	t.Tid = int64(tid)
	t.dbAuth = dbAuth
	t.dbDsn = dbDsn
	t.dbName = dbName

	//make db comon value and check error
	err = t.getDbCache()
	if err != nil {
		fmt.Println("[error]get pid list error", err, t.Tid)
		return nil
	}

	err = t.getPids()
	if err != nil {
		fmt.Println("[error]get pid list error", err, t.Tid)
		t.con.Close()
		return nil
	}

	//if have no pid,no task
	if len(t.pids) == 0 {
		fmt.Println("[notice]this tid have no pid,so pass", t.Tid)
		t.con.Close()
		return nil
	}

	return t
}

func (t *Task) getDbCache() error {
	con, err := sql.Open("mysql", t.dbAuth+"@tcp("+t.dbDsn+")/"+t.dbName+"?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
		return errors.New("connect db error")

	}
	// Third, we wrap in a prepared statement cache for better performance.
	cache := squirrel.NewStmtCacheProxy(con)
	t.dbCache = cache
	t.con = con
	return nil
}

func (t *Task) Do() error {
	pidlen := len(t.pids)
	for i := 0; i < pidlen; i++ {
		t.handlePid(t.pids[i])
	}
	return nil
}

func (t *Task) handlePid(pid int64) error {
	if pid <= 0 {
		fmt.Println("[error]pid <0")
		return errors.New("pid<0")
	}
	post := post.NewPost(t.dbCache, "mysql", pid, t.Tid, false)
	exist := post.PidExists()
	if exist {
		fmt.Println("[info]doing for this pid", post.Message, pid)
		result := post.MoveToSplit()
		if !result {
			fmt.Println("[error] change pid error", pid)
			return errors.New("change pid error")
		}
	} else {
		fmt.Println("[info]not exisit", pid)
		return nil
	}

	return nil

}

func (t *Task) getPids() error {
	//connect db
	db, err := sql.Open("mysql", t.dbAuth+"@tcp("+t.dbDsn+")/"+t.dbName+"?charset=utf8mb4")
	if err != nil {
		fmt.Println("[error] connect db err")

		//fmt.Println("[error] connect db err")
		return errors.New("connect db error")
	}
	defer db.Close()

	//do query part
	rows, err := db.Query("select pid from pre_forum_post where tid=" + strconv.Itoa(int(t.Tid)))
	if err != nil {
		fmt.Println("[error] get pids by tid error: ", t.Tid, err)
		return errors.New("get pids error")
	}

	//get what we need

	pids := make([]int64, 0, 100)
	for rows.Next() {
		var pid int64
		if err := rows.Scan(&pid); err != nil {
			fmt.Println("[error] get pid error after row.next ", err)
			return errors.New("get pid error after row.next")
		}
		fmt.Println("[info]find pid:", pid, "in tid:", t.Tid)
		pids = append(pids, pid)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("[error] rows error agrain ", err)
		return errors.New("rows error agrain")
	}
	t.pids = pids
	return nil
}
