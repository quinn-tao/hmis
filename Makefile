VERSION = 0.0.1
NAME = HMIS

BLD = ./bld

version:
	@echo $(VERSION) 
	@go version

run:
	go run main.go

all: 
	@mkdir -p $(BLD)	  
	go build -o $(BLD) ./...

install-deps:
	go get -u github.com/spf13/cobra@latest

install-devdeps: install-deps
	go install github.com/spf13/cobra-cli@latest

.PHONY: all version deps devdeps
