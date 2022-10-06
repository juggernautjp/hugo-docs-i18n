# Go パラメータ
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMODTIDY=$(GOCMD) mod tidy
BINARY_NAME=hugo-docs-i18n.exe
CONFIG_NAME=hugo-docs-i18n.yaml
TEST_DIR=content/ja
BINARY_UNIX=$(BINARY_NAME)_unix
VER="v0.1.5"
VERM="dev version 0.1.5"

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v -buildvcs=false
test:
		$(GOTEST) -v ./...
test-doci18n:
		$(GOTEST) -v ./doci18n
test-locale:
		$(GOTEST) -v ./locale
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
		rm -f $(CONFIG_NAME)
		rm -rf $(TEST_DIR)
run: 
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOGET) github.com/spf13/cobra
		$(GOGET) github.com/spf13/viper
		$(GOGET) github.com/spf13/pflag
		$(GOMODTIDY)
 git-tags:
		git tag -a $(VER) -m $(VERM)
		git push origin --tags
 git-fix:
		git config --global --add safe.directory F:/Docs.repo/hugo-docs-i18n
