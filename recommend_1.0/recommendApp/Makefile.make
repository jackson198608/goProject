install:
	go install github.com/jackson198608/goProject/recommend_1.0/recommendApp 

test:
	go test github.com/jackson198608/goProject/recommend_1.0/recommendApp -v

build:
	GOOS=linux GOARCH=amd64 go build github.com/jackson198608/goProject/recommend_1.0/recommendApp 


