MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
OUTPUT_DIR := $(MKFILE_DIR)output

build-all:
	if [ ! -d $(OUTPUT_DIR) ]; then mkdir $(OUTPUT_DIR); else rm -Rf $(OUTPUT_DIR)/*; fi
	go mod download
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client_windows_x64.exe client.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o $(OUTPUT_DIR)/dyip-client_windows_x86.exe client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client_linux_x64 client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $(OUTPUT_DIR)/dyip-client_linux_x86 client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-client_linux_arm64 client.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client_darwin_x64 client.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-client_darwin_arm64 client.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server_windows_x64.exe server.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o $(OUTPUT_DIR)/dyip-server_windows_x86.exe server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server_linux_x64 server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $(OUTPUT_DIR)/dyip-server_linux_x86 server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-server_linux_arm64 server.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server_darwin_x64 server.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-server_darwin_arm64 server.go
build:
	if [ ! -d $(OUTPUT_DIR) ]; then mkdir $(OUTPUT_DIR); else rm -Rf $(OUTPUT_DIR)/*; fi
	go mod download
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server server.go
docker:
	docker pull ubuntu:focal
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-client:$(VERSION)-ubuntu -f Dockerfile-Client .; fi
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-client:$(VERSION) -f Dockerfile-Client .; fi
	docker build -t zliea/dyip-client:ubuntu -f Dockerfile-Client .
	docker build -t zliea/dyip-client:latest -f Dockerfile-Client .
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-server:$(VERSION)-ubuntu -f Dockerfile-Server .; fi
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-server:$(VERSION) -f Dockerfile-Server .; fi
	docker build -t zliea/dyip-server:ubuntu -f Dockerfile-Server .
	docker build -t zliea/dyip-server:latest -f Dockerfile-Server .
docker-alpine:
	docker pull alpine:latest
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
version:
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_windows_x64.exe $(OUTPUT_DIR)/dyip-client_$(VERSION)_windows_x64.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_windows_x86.exe $(OUTPUT_DIR)/dyip-client_$(VERSION)_windows_x86.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_linux_x64 $(OUTPUT_DIR)/dyip-client_$(VERSION)_linux_x64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_linux_x86 $(OUTPUT_DIR)/dyip-client_$(VERSION)_linux_x86; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_linux_arm64 $(OUTPUT_DIR)/dyip-client_$(VERSION)_linux_arm64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_darwin_x64 $(OUTPUT_DIR)/dyip-client_$(VERSION)_darwin_x64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_darwin_arm64 $(OUTPUT_DIR)/dyip-client_$(VERSION)_darwin_arm64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_windows_x64.exe $(OUTPUT_DIR)/dyip-server_$(VERSION)_windows_x64.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_windows_x86.exe $(OUTPUT_DIR)/dyip-server_$(VERSION)_windows_x86.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_linux_x64 $(OUTPUT_DIR)/dyip-server_$(VERSION)_linux_x64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_linux_x86 $(OUTPUT_DIR)/dyip-server_$(VERSION)_linux_x86; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_linux_arm64 $(OUTPUT_DIR)/dyip-server_$(VERSION)_linux_arm64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_darwin_x64 $(OUTPUT_DIR)/dyip-server_$(VERSION)_darwin_x64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_darwin_arm64 $(OUTPUT_DIR)/dyip-server_$(VERSION)_darwin_arm64; fi
all: clean build-all
release: clean build-all version
image: clean build docker
image-push: image push
image-alpine: clean build docker-alpine
image-alpine-push: image-alpine push-alpine