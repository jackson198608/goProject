package tools

import "github.com/go-xorm/xorm"

func GetMysqlConnect(mysqlInfo []string) ([]*xorm.Engine, error) {

	//if you do not need mysql for job func
	if mysqlInfo == nil {
		return nil, nil
	}

	// if you need make it for you, and info must be correct
	mysqls := []*xorm.Engine{}
	for _, mysqlInfo := range mysqlInfo {
		x, err := GetMysqlSingleConnect(mysqlInfo)
		if err != nil {
			//close former connection
			closeMysqlConn(mysqls)
			return nil, err
		}
		mysqls = append(mysqls, x)
	}
	return mysqls, nil

}
func GetMysqlSingleConnect(mysqlInfo string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", mysqlInfo)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

func closeMysqlConn(mysqlConns []*xorm.Engine) {
	if mysqlConns == nil {
		return
	}
	for _, conn := range mysqlConns {
		conn.Close()
	}
	return
}
