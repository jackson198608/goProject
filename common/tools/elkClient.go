package tools

import (
	"github.com/olivere/elastic"
	"log"
	"time"
)

func NewClient(elkDsn string) (*Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://".elkDsn),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetRetrier(NewCustomRetrier()),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	return client, err
}
