package tools

import ()

//return the dns string for mysql
func GetMysqlDsn(dbAuth string, dbDsn string, dbName string) string {
	return dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
}
