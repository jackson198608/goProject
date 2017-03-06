package post

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"testing"
)

func getDbCache() squirrel.DBProxyBeginner {
	dbName := "new_dog123"
	con, err := sql.Open("mysql", "root:my-secret-pw@tcp(127.0.0.1:3306)/"+dbName+"?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
		return nil
	}
	// Third, we wrap in a prepared statement cache for better performance.
	cache := squirrel.NewStmtCacheProxy(con)
	return cache
}

func NewWithTidPid() *Post {
	cache := getDbCache()
	post := NewPost(cache, "mysql", 65441666, 4411466, false)
	exist := post.PidExists()
	if exist {
		fmt.Println(post.Message)
	} else {

		fmt.Println("not exisit")
	}
	return post
}

func TestNewWithTidPid(t *testing.T) {
	NewWithTidPid()
}

func TestMove(t *testing.T) {
	p := NewWithTidPid()
	p.MoveToSplit()
}

func TestInsert(t *testing.T) {
	// Boilerplate DB setup.
	// First, we need to know the database driver.
	// Second, we need a database connection.
	// Create an empty new user and give it some properties.
	cache := getDbCache()
	post := NewPost(cache, "mysql", 0, 0, false)
	post.Pid = 65441670
	post.Fid = 38
	post.Tid = 4411466
	post.First = false
	post.Author = "gm_DZVOIoz4t66E"
	post.Authorid = 1850916
	post.Subject = ""
	post.Dateline = 1488202168
	post.Message = "囡囡小白来报道"
	post.Useip = "222.64.32.123"
	post.Invisible = false
	post.Anonymous = false
	post.Usesig = true
	post.Htmlon = false
	post.Bbcodeoff = false
	post.Smileyoff = false
	post.Parseurloff = false
	post.Attachment = false
	post.Rate = 0
	post.Ratetimes = 0
	post.Status = 136
	post.Tags = ""
	post.Comment = false
	post.Replycredit = 0

	// Insert this as a new record.
	if err := post.Insert(); err != nil {
		panic(err.Error())
	}

}
