package fileEngine

import (
	"testing"
)

const filePath = "/data/targetRefer"

func TestFileEngine(t *testing.T) {
	f, _ := NewFileEngine(filePath)
	for {
		line, err := f.ReadLine()
		if err != nil {
			break
		}
		t.Log(line)
	}

	f.Close()

}
