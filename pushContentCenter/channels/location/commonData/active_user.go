package commonData

import (
	// "errors"
	// "fmt"
	mgo "gopkg.in/mgo.v2"
	// "math"
	// "reflect"
	// "strconv"
)

func LoadDataToHashmap(mc *mgo.Session) map[int]bool {
	var m map[int]bool
	m = make(map[int]bool)

	var uids []int
	c := mc.DB("ActiveUser").C("active_user")
	err := c.Find(nil).Distinct("uid", &uids)
	if err != nil {
		panic(err)
		return m
	}
	for i := 0; i < len(uids); i++ {
		m[uids[i]] = true
	}
	return m
}
