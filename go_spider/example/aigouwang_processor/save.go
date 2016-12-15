package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/hu17889/go_spider/core/common/page"
	"os"
	"path"
)

func checkCanSave(p *page.Page) bool {
	return true
}

func saveImage(p *page.Page) bool {
	//judge if status is 200
	if !checkCanSave(p) {
		return false
	}

	//get md5
	h := md5.New()
	url := p.GetRequest().Url
	h.Write([]byte(url))
	md5Str := hex.EncodeToString(h.Sum(nil))

	//get fullpath
	abPath := getPath(md5Str)
	fullDirPath := saveDir + abPath
	err := os.MkdirAll(fullDirPath, 0664)
	if err != nil {
		logger.Println("[error]create dir error:", err, " ", fullDirPath, " ", url)
		return false
	}

	//save file

	fileName := fullDirPath + "/" + path.Base(url)
	result, err := os.Create(fileName)
	if err != nil {
		logger.Println("[error]create file error:", err, " ", fileName, " ", url)
		return false
	}
	logger.Println("[info] save image:", url)
	logger.Println("[info] save in:", fileName)
	logger.Println("[info] save len:", len(p.GetBodyStr()))
	result.WriteString(p.GetBodyStr())
	result.Close()

	return true
}

func save(p *page.Page) bool {
	//judge if status is 200
	if !checkCanSave(p) {
		return false
	}

	//get md5
	h := md5.New()
	url := p.GetRequest().Url
	h.Write([]byte(url))
	md5Str := hex.EncodeToString(h.Sum(nil))

	//get fullpath
	abPath := getPath(md5Str)
	fullDirPath := saveDir + abPath
	err := os.MkdirAll(fullDirPath, 0664)
	if err != nil {
		logger.Println("[error]create dir error:", err, " ", fullDirPath, " ", url)
		return false
	}

	//save file
	fileName := fullDirPath + "/" + path.Base(url)
	result, err := os.Create(fileName)
	if err != nil {
		logger.Println("[error]create file error:", err, " ", fileName, " ", url)
		return false
	}
	logger.Println("[info] save page:", url)
	logger.Println("[info] save in:", fileName)
	result.WriteString(url + "\n")
	result.WriteString(p.GetBodyStr())
	result.Close()

	return true
}

func getPath(md5Str string) string {
	abPath := make([]byte, 48, 48)
	j := 0
	p := 0
	for i := 0; i <= 16; i++ {
		if p >= 50 || j >= 32 {
			break
		}
		abPath[p] = byte('/')
		abPath[p+1] = md5Str[j]
		abPath[p+2] = md5Str[j+1]
		p = p + 3
		j = j + 2
	}
	return string(abPath)
}
