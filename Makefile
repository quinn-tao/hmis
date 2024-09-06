NAME = HMIS
BLD = ./bld

.fmt:
	go fmt 

.install-deps:
	@echo "Installing dependencies ..."
	go get -u github.com/spf13/cobra@latest     # cli parser
	go get -u github.com/spf13/viper@latest     # cfg parser 
	go get -u github.com/jedib0t/go-pretty/v6   # pretty printer
	go get github.com/jedib0t/go-pretty/v6/text@v6.5.9
	go get -u gopkg.in/yaml.v3					# yaml parser
	go get -u github.com/mattn/go-sqlite3

.install-devdeps: .install-deps
	@echo "Installing dev-dependencies ..."
	go install github.com/spf13/cobra-cli@latest 

all: .install-devdeps .fmt 
	@mkdir -p $(BLD)	  
	@echo "Building projects ..."
	go build -o $(BLD)/hmis main.go

build:
	go build -o $(BLD)/hmis main.go

run:
	go run main.go

clean:
	go clean 
	rm $(BLD)/*

.PHONY: all clean .version .install-deps .intall-devdeps \
	.fmt run build
