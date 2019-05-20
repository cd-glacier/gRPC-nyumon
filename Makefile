.PHONY: help
.DEFAULT_GOAL := help

echo-protoc: ## generate proto code
	protoc --proto_path=./echo/proto --go_out=plugins=grpc:./echo/proto/ ./echo/proto/echo.proto

echo-server: ## start echo grpc server
	go run ./echo/server/

help: ## show help to make
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
