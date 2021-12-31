MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
OUTPUT_DIR := $(MKFILE_DIR)output

build:
	if [ ! -d $(OUTPUT_DIR) ]; then mkdir $(OUTPUT_DIR); else rm -Rf $(OUTPUT_DIR)/*; fi
	go mod download
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server server.go
docker:
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-client:$(VERSION)-ubuntu -f Dockerfile-Client .; fi
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-client:$(VERSION) -f Dockerfile-Client .; fi
	docker build -t zliea/dyip-client:ubuntu -f Dockerfile-Client .
	docker build -t zliea/dyip-client:latest -f Dockerfile-Client .
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-server:$(VERSION)-ubuntu -f Dockerfile-Server .; fi
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-server:$(VERSION) -f Dockerfile-Server .; fi
	docker build -t zliea/dyip-server:ubuntu -f Dockerfile-Server .
	docker build -t zliea/dyip-server:latest -f Dockerfile-Server .
docker-alpine:
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-client:$(VERSION)-alpine -f Dockerfile-Client-Alpine .; fi
	docker build -t zliea/dyip-client:alpine -f Dockerfile-Client-Alpine .
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-server:$(VERSION)-alpine -f Dockerfile-Server-Alpine .; fi
	docker build -t zliea/dyip-server:alpine -f Dockerfile-Server-Alpine .
push:
	if [ -n "$(VERSION)" ]; then docker push zliea/dyip-client:$(VERSION)-ubuntu; fi
	if [ -n "$(VERSION)" ]; then docker push zliea/dyip-client:$(VERSION); fi
	docker push zliea/dyip-client:ubuntu
	docker push zliea/dyip-client:latest
	if [ -n "$(VERSION)" ]; then docker push zliea/dyip-server:$(VERSION)-ubuntu; fi
	if [ -n "$(VERSION)" ]; then docker push zliea/dyip-server:$(VERSION); fi
	docker push zliea/dyip-server:ubuntu
	docker push zliea/dyip-server:latest
push-alpine:
	if [ -n "$(VERSION)" ]; then docker push zliea/dyip-client:$(VERSION)-alpine; fi
	docker push zliea/dyip-client:alpine
	if [ -n "$(VERSION)" ]; then docker push zliea/dyip-server:$(VERSION)-alpine; fi
	docker push zliea/dyip-server:alpine
clean:
	rm -Rf $(OUTPUT_DIR)
all: clean build docker