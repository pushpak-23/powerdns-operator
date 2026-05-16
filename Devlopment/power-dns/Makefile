GO ?= go
IMG ?= powerdns-platform-operator:dev

.PHONY: test build docker-build tidy

test:
	PATH=/usr/local/go/bin:$(PATH) $(GO) test ./...

build:
	PATH=/usr/local/go/bin:$(PATH) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o bin/manager ./cmd/manager

docker-build:
	docker build -t $(IMG) .

tidy:
	PATH=/usr/local/go/bin:$(PATH) $(GO) mod tidy