package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
)

type InfoAsk struct {
	db           *sql.DB
	id           int
	session      *mgo.Session
	templateType string
	templatefile string
	saveDir      string
	tidStart     string
	tidEnd       string
	domain       string
}

func NewAskInfo(logLevel int, id int, db *sql.DB, session *mgo.Session, taskNewArgs []string) *InfoAsk {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(InfoAsk)
	e.db = db
	e.id = id
	e.session = session
	e.templateType = taskNewArgs[3]
	e.templatefile = taskNewArgs[4]
	e.saveDir = taskNewArgs[5]
	e.tidStart = taskNewArgs[6]
	e.tidEnd = taskNewArgs[7]
	e.domain = taskNewArgs[8]
	return e
}

func (e *InfoAsk) CreateAskHtmlContent(id int) error {
	question := LoadQuestionById(id, e.db)
	if question.Id <= 0 || question == nil {
		logger.Info("ask_question is not exist id=", id)
		return nil
	}
	fmt.Println(question.Images)
	//相关帖子 eg:tid=12
	// relateThread := e.relateThread(id, thread.Fid, e.db, e.session)
	// //相关问答 eg:tid=12
	// relateAsk := relateAsk(id, e.db, e.session)
	// if relateAsk == "" {
	// 	relateAsk = relateDefaultAsk
	// }
	// //相关犬种 eg:tid=4682521
	// relateDogs := relateDogs(id, e.db, e.session, e.templateType)
	// posts := LoadPostsByTid(id, thread.Posttableid, e.db)
	// if posts == nil {
	// 	logger.Info("post is not exist tid=", tid)
	// 	return nil
	// }
	return nil
}

func DefaultDoctors(e.db) string {
	doctors := LoadHealthDoctor(db)
	if doctors == nil {
		return ""
	}
	var s string = ""
	for _, v := range doctors {
		s += "<dl><a href=\"http://a.app.qq.com/o/simple.jsp?pkgname=com.goumin.forum\" class=\"doctor-avatar ui-link\"><dt><mip-img src=\"" + v.Avatar + "\"><em>" + v.Name + "</em></dt><dd>" + v.Hospital + "</dd></a></dl>"
	}
	return s
}
