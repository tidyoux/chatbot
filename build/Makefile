
macos:
	@go install ../cmd
	go build -o ./bin/chatbot-macos ../cmd

linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/chatbot-linux ../cmd
	
.PHONY: macos linux