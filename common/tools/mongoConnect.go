package tools

import (
	mgo "gopkg.in/mgo.v2"
)

func GetReplicaConnecting(mgoInfo []string) (*mgo.Session, error) {
	const (
		Username       = ""
		Password       = ""
		Database       = ""
		ReplicaSetName = "goumin"
	)
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          mgoInfo,
		Username:       Username,
		Password:       Password,
		Database:       Database,
		ReplicaSetName: ReplicaSetName,
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func GetStandAloneConnecting(mgoInfo string) (*mgo.Session, error) {
	session, err := mgo.Dial(mgoInfo)
	if err != nil {
		return nil, err
	}

	return session, nil
}
