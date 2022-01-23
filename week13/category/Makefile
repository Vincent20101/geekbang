GOPATH:=$(shell go env GOPATH)
.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/category/category.proto

.PHONY: build
build: 

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o category-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t category-service:latest
