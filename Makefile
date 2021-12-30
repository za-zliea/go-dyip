MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
OUTPUT_DIR := $(MKFILE_DIR)output

build:
	if [ ! -d $(OUTPUT_DIR) ]; then mkdir $(OUTPUT_DIR); else rm -Rf $(OUTPUT_DIR)/*; fi
	go mod download
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server server.go
docker:
	docker build -t dyip-client:ubuntu -f Dockerfile-Client .
	docker build -t dyip-client:latest -f Dockerfile-Client .
	docker build -t dyip-server:ubuntu -f Dockerfile-Server .
	docker build -t dyip-server:latest -f Dockerfile-Server .
docker-alpine:
	docker build -t dyip-client:alpine -f Dockerfile-Client-Alpine .
	docker build -t dyip-server:alpine -f Dockerfile-Server-Alpine .
clean:
	rm -Rf $(OUTPUT_DIR)
all: clean build docker
