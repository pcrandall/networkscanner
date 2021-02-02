currentDir = $(shell pwd)

## installation
install:
	@echo "==> installing go dependencies"
	@go mod download
.PHONY: install

run:
	@echo "==> running network scanner"
	@go run *.go
.PHONY: run

buildwindows:
	@echo "==> building network scanner for windows"
	@${currentDir}/scripts/buildwin.sh
.PHONY: build

buildlinux:
	@echo "==> building network scanner for linux"
	@${currentDir}/scripts/buildlinux.sh
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
