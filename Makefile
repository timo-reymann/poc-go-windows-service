.PHONY: help
SHELL := /bin/bash

help: ## Display this help page
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

build-binary: ## Build binaries for windows
	@GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64.exe
	@GOOS=windows GOARCH=arm go build -o dist/windows-arm64.exe

register-service: ## Register the service using sc
	sc.exe create poc-go-windows-service start= auto binPath=$(shell pwd)/dist/windows-$(shell arch).exe
