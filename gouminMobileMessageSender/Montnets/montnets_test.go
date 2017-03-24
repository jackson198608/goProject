package Montnets

import (
	"testing"
	//"github.com/donnie4w/go-logger/logger"
)

func TestSend(t *testing.T) {
	m := NewMontnets(0, "18210091845", "1319")
	if m != nil {
		m.send()
	}
}
