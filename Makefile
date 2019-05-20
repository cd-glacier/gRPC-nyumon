.PHONY: help server echo chat
.DEFAULT_GOAL := help

protoc: ## generate proto code
	protoc --proto_path=./echo --go_out=plugins=grpc:./echo/ ./echo/echo.proto
	protoc --proto_path=./downloader --go_out=plugins=grpc:./downloader/ ./downloader/downloader.proto
	protoc --proto_path=./uploader --go_out=plugins=grpc:./uploader/ ./uploader/uploader.proto
	protoc --proto_path=./chat --go_out=plugins=grpc:./chat/ ./chat/chat.proto

server: ## start file grpc server
	go run ./server/

echo: ## echo "hello"
	go run ./client/ --mode=echo --message="hello"

file-download: ## download "resource.txt"
	go run ./client/ --mode=download --name=resource.txt

file-upload: ## upload "resource.txt"
	go run ./client/ --mode=upload --name=resource.txt

chat: ## chat
	go run ./client/ --mode=chat --name="hyoga"

help: ## show help to make
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
