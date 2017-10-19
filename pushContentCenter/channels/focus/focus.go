package focus

import (
	"errors"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
)

type Focus struct {
	mysqlXorm *xorm.Engine
	mongoConn *mgo.Session
	jobstr    string
	jsonData  *jsonColumn
}

type jsonColumn struct {
	//json column
}

func NewFocus(mysqlXorm *xorm.Engine, mongoConn *mgo.Session, jobStr string) *Focus {
	//@todo checkparams
	f := new(Focus)
	if f == nil {
		return nil
	}

	//@todo pass params
	err := f.parseJson()
	if err == nil {
		return nil
	}

	return f

}

func (f *Focus) Do() error {
	page := f.getPersionsPageNum()
	if page <= 0 {
		return nil
	}

	for i := 1; i <= page; i++ {
		currentPersionList := f.getPersons(page)
		f.pushPersons(currentPersionList)

	}

}

//change json colum to object private member
func (f *Focus) parseJson() error {

}

func (f *Focus) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("you have no person to push " + f.jobstr)
	}

	for _, person := range persons {
		err := f.pushPerson(person)
		if err != nil {
			//@todo if err times < 5 ,just print log
			//      if err times > 5 ,return err
		}
	}
	return nil

}

func (f *Focus) pushPerson(person int) error {

	return nil
}

//@todo how to remove duplicate uid from to lists
func (f *Focus) getPersons(page int) []int {

	return 0
}

func (f *Focus) getPersionsPageNum() int {

	return 0
}
