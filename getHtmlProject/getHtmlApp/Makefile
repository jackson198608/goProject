install:
	go install github.com/jackson198608/goProject/getHtmlProject/getHtmlApp
	cp data/config* /etc/
	cp ../../../../../bin/getHtmlApp /usr/local/bin/
runAsk:
	getHtmlApp asksave /etc/configask.json
	getHtmlApp asksave /etc/configask.json update
	getHtmlApp ask /etc/configask.json
runThread:
	getHtmlApp threadsave /etc/configthread.json update
	getHtmlApp thread /etc/configthread.json
