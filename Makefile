.PHONY: help
.DEFAULT_GOAL := help

echo-protoc: ## generate proto code
	protoc --proto_path=./echo/proto --go_out=plugins=grpc:./echo/proto/ ./echo/proto/echo.proto

echo-server: ## start echo grpc server
	go run ./echo/server/

echo-client: ## start echo client and send "test"
	go run ./echo/client/ test

file-protoc: ## generate proto code
	protoc --proto_path=./file/proto --go_out=plugins=grpc:./file/proto/ ./file/proto/downloader.proto
	protoc --proto_path=./file/proto --go_out=plugins=grpc:./file/proto/ ./file/proto/uploader.proto

file-server: ## start file grpc server
	go run ./file/server/

file-client: ## start file client and get "resource.txt"
	go run ./file/client/ resource.txt

help: ## show help to make
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
