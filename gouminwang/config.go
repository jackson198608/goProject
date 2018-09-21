package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn  string
	dbName string
	dbAuth string
}

func loadConfig(args []string) {
	var config *jconfig.Config

	if len(args) >= 3 {
		config = jconfig.LoadConfig(args[2])
	} else {
		config = jconfig.LoadConfig("/etc/gouminwang.json")
	}
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
}
