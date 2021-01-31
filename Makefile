currentDir = $(shell pwd)

## installation
install:
	@echo "==> installing go dependencies"
	@go mod download
.PHONY: install

run:
	@echo "==> running network scanner"
	@go run .
.PHONY: run

buildwin:
	@echo "==> building network scanner for windows"
	@GOOS=windows GOARCH=386 go build .
.PHONY: build

build:
	@echo "==> building network scanner"
	@go build .
.PHONY: build

git:
	@echo "==> adding git tracked files"
	@git add -u
	@git commit
	@echo "==> pushing to git remote"
	@git push origin
.PHONY: git

clean:
	@go clean --cache
	@go mod tidy
	@git clean -f
.PHONY: clean
