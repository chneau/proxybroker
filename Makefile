.SILENT:
.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.PHONY: run exec build clean deps

name=$(shell basename $(CURDIR))

run: build exec clean

exec:
	./bin/${name}

build:
	CGO_ENABLED=0 go build -trimpath -o bin/${name} -ldflags '-s -w -extldflags "-static"'

clean:
	rm -rf bin

test:
	go test -cover -count=1 ./...



deps:
	rm -rf go.mod go.sum
	go mod init || true
	go mod tidy
	go mod verify
