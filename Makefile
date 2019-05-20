.PHONY: help
.DEFAULT_GOAL := help

echo-protoc: ## generate proto code
	protoc --proto_path=./echo/proto --go_out=plugins=grpc:./echo/proto/ ./echo/proto/echo.proto

echo-server: ## start echo grpc server
	go run ./echo/server/

echo-client: ## start echo client and send "test"
	go run ./echo/client/ test

file-protoc: ## generate proto code
	protoc --proto_path=./downloader --go_out=plugins=grpc:./downloader/ ./downloader/downloader.proto
	protoc --proto_path=./uploader --go_out=plugins=grpc:./uploader/ ./uploader/uploader.proto

file-server: ## start file grpc server
	go run ./server/

file-download: ## download "resource.txt"
	go run ./client/ --mode=download --filename=resource.txt

file-upload: ## upload "resource.txt"
	go run ./client/ --mode=upload --filename=resource.txt

help: ## show help to make
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
