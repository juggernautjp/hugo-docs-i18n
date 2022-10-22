SHELL=/bin/bash
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
DIST_DIR=dist
TEST_VER=docsi18n/version.json
BUG_DIR=content/en/myshowcase
BINARY_UNIX=$(BINARY_NAME)_unix
# GIT_VER=$$(./hugo-docs-i18n.exe version -g)
GIT_VER=`./hugo-docs-i18n.exe version -g`
# VER_MSG=$$(./hugo-docs-i18n.exe version -m)
VER_MSG=`./hugo-docs-i18n.exe version -m`
# NOW=`date`
NOW=$$(date)

all: test build
build: $(BINARY_NAME)

$(BINARY_NAME):	doci18n/* locale/* cmd/*
		$(GOBUILD) -o $(BINARY_NAME) -v -buildvcs=false
test: 
		$(GOTEST) -v ./...
test-doci18n:
		$(GOTEST) -v ./doci18n
test-locale:
		$(GOTEST) -v ./locale
test-cmd:
		$(GOTEST) -v ./cmd
clean:
		$(GOCLEAN)
		rm -f $(TEST_VER)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
		rm -f $(CONFIG_NAME)
		rm -rf $(TEST_DIR)
		rm -rf $(DIST_DIR)
run: $(BINARY_NAME)
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOGET) github.com/spf13/cobra
		$(GOGET) github.com/spf13/viper
		$(GOGET) github.com/spf13/pflag
		$(GOMODTIDY)
git-tags: $(BINARY_NAME)
		git tag -a $(GIT_VER) -m "$(VER_MSG)"
		git tag -l
		git ls-remote https://github.com/juggernautjp/hugo-docs-i18n/
		echo 'execute "git push origin --tags"'
#		git push origin --tags
git-fix:
		git config --global --add safe.directory F:/Docs.repo/hugo-docs-i18n
now: $(BINARY_NAME)
		@echo $(NOW) $(GIT_VER) "$(VER_MSG)"
