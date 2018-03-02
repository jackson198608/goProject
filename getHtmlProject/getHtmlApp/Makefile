install:
	go install github.com/jackson198608/goProject/getHtmlProject/getHtmlApp
runAsk:
	getHtmlApp asksave /etc/configask.json
	getHtmlApp asksave /etc/configask.json update
	getHtmlApp ask /etc/configask.json
runThread:
	getHtmlApp threadsave /etc/configthread.json update
	getHtmlApp thread /etc/configthread.json
runForum:
	getHtmlApp bbsindexsave /etc/configforum.json　
	getHtmlApp bbsindex /etc/configforum.json　
	getHtmlApp forumsave /etc/configforum.json
	getHtmlApp forum /etc/configforum.json

build:
	GOOS=linux GOARCH=amd64 go build github.com/jackson198608/goProject/getHtmlProject/getHtmlApp
