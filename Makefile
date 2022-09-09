# Go パラメータ
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMODTIDY=$(GOCMD) mod tidy
BINARY_NAME=hugo-docs-i18n.exe
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run: 
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOGET) github.com/spf13/cobra
		$(GOGET) github.com/spf13/viper
		$(GOGET) github.com/spf13/pflag
		$(GOMODTIDY)

