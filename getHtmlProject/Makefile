install:
	go install github.com/jackson198608/goProject/getHtml
	cp data/config* /etc/
	cp ../../../../../bin/getHtml /usr/local/bin/
runAsk:
	getHtml asksave /etc/configask.json
	getHtml ask /etc/configask.json
runThread:
	getHtml threadsave /etc/configthread.json
	getHtml thread /etc/configthread.json
build:
	GOOS=linux GOARCH=amd64 go build github.com/jackson198608/goProject/getHtmlProject 
	
