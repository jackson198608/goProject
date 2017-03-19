package post

import (
	"database/sql"
	"fmt"
	"github.com/jackson198608/squirrel"
	"testing"
)

func getDbCache() squirrel.DBProxyBeginner {
	dbName := "new_dog123"
	con, err := sql.Open("mysql", "dog123:dog123@tcp(210.14.154.198:3306)/"+dbName+"?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
		return nil
	}
	// Third, we wrap in a prepared statement cache for better performance.
	cache := squirrel.NewStmtCacheProxy(con)
	return cache
}

func NewWithTidPid(pid int64, tid int64) (*Post, bool) {
	cache := getDbCache()
	post := NewPost(0, cache, "mysql", pid, tid, false)
	exist := post.PidExists()
	if exist {
		fmt.Println("exist")
		fmt.Println(post.Message)
	} else {

		fmt.Println("not exisit")
		return post, false
	}
	return post, true
}

func TestNewWithTidPid(t *testing.T) {
	fmt.Println("testing new tid pid")
	NewWithTidPid(1, 1)
}

func TestMove(t *testing.T) {
	p, isExist := NewWithTidPid(18994058, 12322)
	if isExist {
		p.MoveToSplit()
	}

}

func croutineJob(c chan int, pid int64, tid int64) {
	p, isExist := NewWithTidPid(pid, tid)
	if isExist {
		p.MoveToSplit()
	}

	c <- 1
}

func TestMulti(t *testing.T) {
	tids := [10]int64{14716, 16567, 16741, 16540, 14716, 9223, 16437, 15326, 724, 16679}
	pids := [10]int64{113344, 113345, 113346, 113347, 113348, 113349, 113350, 113351, 113352, 113354}

	taskNum := 10
	c := make(chan int, taskNum)
	for i := 0; i < taskNum; i++ {
		go croutineJob(c, pids[i], tids[i])
	}

	for i := 0; i < taskNum; i++ {
		<-c
	}

}

func TestInsert(t *testing.T) {
	// Boilerplate DB setup.
	// First, we need to know the database driver.
	// Second, we need a database connection.
	// Create an empty new user and give it some properties.
	cache := getDbCache()
	post := NewPost(0, cache, "mysql", 65441670, 4411466, true)
	post.Pid = 65441670
	post.Fid = 38
	post.Tid = 4411466
	post.First = 0
	post.Author = "gm_DZVOIoz4t66E"
	post.Authorid = 1850916
	post.Subject = ""
	post.Dateline = 1488202168
	post.Message = "囡囡小白来报道"
	post.Useip = "222.64.32.123"
	post.Invisible = 0
	post.Anonymous = 0
	post.Usesig = 1
	post.Htmlon = 0
	post.Bbcodeoff = 0
	post.Smileyoff = 0
	post.Parseurloff = 0
	post.Attachment = 0
	post.Rate = 0
	post.Ratetimes = 0
	post.Status = 136
	post.Tags = ""
	post.Comment = 0
	post.Replycredit = 0

	// Insert this as a new record.
	if err := post.Insert(); err != nil {
		fmt.Println("insert sql error", err)
	}

}
