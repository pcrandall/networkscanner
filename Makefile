currentDir = $(shell pwd)


## installation
install:
	@echo "==> installing go dependencies"
	@go mod download
.PHONY: install

# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif
run:
	@echo "==> running network scanner"
	@go run . $(RUN_ARGS)
.PHONY: run

buildwin:
	@echo "==> building network scanner for windows"
	@GOOS=windows GOARCH=386 go build -tags windows .
.PHONY: build

build:
	@echo "==> building network scanner"
	@go build -tags linux .
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
