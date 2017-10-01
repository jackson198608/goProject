package readFileEngine

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type FileEngine struct {
	filePath   string
	fileFp     *os.File
	fileReader *bufio.Reader
}

func NewFileEngine(filePath string) (*FileEngine, error) {
	fileEngine := new(FileEngine)
	if fileEngine == nil {
		return nil, errors.New("create struct error")
	}
	fileEngine.filePath = filePath
	err := fileEngine.open()
	if err != nil {
		return nil, err
	}

	return fileEngine, nil
}

func (f *FileEngine) open() error {
	file, err := os.Open(f.filePath)
	if err != nil {
		return err
	}

	rd := bufio.NewReader(f)

	f.fileFp = file
	f.fileReader = rd
	return nil
}

func (f *FileEngine) Close() error {
	err := f.fileFp.Close()
	return err
}

func (f *FileEngine) ReadLine() (string, error) {
	line, err := f.fileReader.ReadString('\n') //以'\n'为结束符读入一行
	if err != nil || io.EOF == err {
		return "", err
	}
	line = strings.Replace(line, "\n", "", -1)
	line = strings.Replace(line, " ", "", -1)
	line1 := line[0 : len(line)-1]
	return string(line1)
}
