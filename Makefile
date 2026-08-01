MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
OUTPUT_DIR := $(MKFILE_DIR)output

.PHONY: all build build-all clean docker docker-alpine frontend-build frontend-clean image image-alpine image-push image-alpine-push push push-alpine release version

frontend-clean:
	@if command -v pnpm >/dev/null 2>&1; then \
    	echo ">> cleaning frontend with pnpm"; \
    	cd src/frontend && pnpm clean; \
    elif command -v npx >/dev/null 2>&1; then \
    	echo ">> pnpm not found, falling back to npx pnpm via corepack"; \
    	cd src/frontend && npx --yes pnpm@9 clean; \
    else \
    	echo "!! pnpm/npx not found — skipping frontend clean (using existing src/frontend/dist if any)"; \
    fi

frontend-build:
	@if command -v pnpm >/dev/null 2>&1; then \
		echo ">> building frontend with pnpm"; \
		cd src/frontend && pnpm install --frozen-lockfile=false && pnpm build; \
	elif command -v npx >/dev/null 2>&1; then \
		echo ">> pnpm not found, falling back to npx pnpm via corepack"; \
		cd src/frontend && npx --yes pnpm@9 install --frozen-lockfile=false && npx --yes pnpm@9 build; \
	else \
		echo "!! pnpm/npx not found — skipping frontend build (using existing src/frontend/dist if any)"; \
	fi

build-all: frontend-build
	if [ ! -d $(OUTPUT_DIR) ]; then mkdir $(OUTPUT_DIR); else rm -Rf $(OUTPUT_DIR)/*; fi
	go mod download
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client_windows_x64.exe src/client.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o $(OUTPUT_DIR)/dyip-client_windows_x86.exe src/client.go
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-client_windows_arm64.exe src/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client_linux_x64 src/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $(OUTPUT_DIR)/dyip-client_linux_x86 src/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-client_linux_arm64 src/client.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client_darwin_x64 src/client.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-client_darwin_arm64 src/client.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server_windows_x64.exe src/server.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o $(OUTPUT_DIR)/dyip-server_windows_x86.exe src/server.go
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-server_windows_arm64.exe src/server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server_linux_x64 src/server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $(OUTPUT_DIR)/dyip-server_linux_x86 src/server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-server_linux_arm64 src/server.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server_darwin_x64 src/server.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/dyip-server_darwin_arm64 src/server.go
build: frontend-build
	if [ ! -d $(OUTPUT_DIR) ]; then mkdir $(OUTPUT_DIR); else rm -Rf $(OUTPUT_DIR)/*; fi
	go mod download
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-client src/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/dyip-server src/server.go
docker:
	docker pull zliea/ubuntu:noble
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-client:$(VERSION)-ubuntu -f docker/Dockerfile-Client .; fi
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-client:$(VERSION) -f docker/Dockerfile-Client .; fi
	docker build -t zliea/dyip-client:ubuntu -f docker/Dockerfile-Client .
	docker build -t zliea/dyip-client:latest -f docker/Dockerfile-Client .
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-server:$(VERSION)-ubuntu -f docker/Dockerfile-Server .; fi
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-server:$(VERSION) -f docker/Dockerfile-Server .; fi
	docker build -t zliea/dyip-server:ubuntu -f docker/Dockerfile-Server .
	docker build -t zliea/dyip-server:latest -f docker/Dockerfile-Server .
docker-alpine:
	docker pull alpine:3
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-client:$(VERSION)-alpine -f docker/Dockerfile-Client-Alpine .; fi
	docker build -t zliea/dyip-client:alpine -f docker/Dockerfile-Client-Alpine .
	if [ -n "$(VERSION)" ]; then docker build -t zliea/dyip-server:$(VERSION)-alpine -f docker/Dockerfile-Server-Alpine .; fi
	docker build -t zliea/dyip-server:alpine -f docker/Dockerfile-Server-Alpine .
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
clean: frontend-clean
	rm -Rf $(OUTPUT_DIR)
	go clean --cache
version:
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_windows_x64.exe $(OUTPUT_DIR)/dyip-client_$(VERSION)_windows_x64.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_windows_x86.exe $(OUTPUT_DIR)/dyip-client_$(VERSION)_windows_x86.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_windows_arm64.exe $(OUTPUT_DIR)/dyip-client_$(VERSION)_windows_arm64.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_linux_x64 $(OUTPUT_DIR)/dyip-client_$(VERSION)_linux_x64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_linux_x86 $(OUTPUT_DIR)/dyip-client_$(VERSION)_linux_x86; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_linux_arm64 $(OUTPUT_DIR)/dyip-client_$(VERSION)_linux_arm64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_darwin_x64 $(OUTPUT_DIR)/dyip-client_$(VERSION)_darwin_x64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-client_darwin_arm64 $(OUTPUT_DIR)/dyip-client_$(VERSION)_darwin_arm64; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_windows_x64.exe $(OUTPUT_DIR)/dyip-server_$(VERSION)_windows_x64.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_windows_x86.exe $(OUTPUT_DIR)/dyip-server_$(VERSION)_windows_x86.exe; fi
	if [ -n "$(VERSION)" ]; then mv $(OUTPUT_DIR)/dyip-server_windows_arm64.exe $(OUTPUT_DIR)/dyip-server_$(VERSION)_windows_arm64.exe; fi
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
