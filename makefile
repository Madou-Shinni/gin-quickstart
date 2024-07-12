VERSION=$(shell git describe --tags --always)

.PHONY: init
# initialize env
init:
	go install github.com/swaggo/swag/cmd/swag@v1.7.9
	go install github.com/Madou-Shinni/gctl@latest

# build
build:
	docker build -t gin-quickstart:$(VERSION) .

# sync api
api-sync:
	swag init && go run cmd/auto/main.go
