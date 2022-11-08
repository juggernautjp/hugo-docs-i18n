SHELL=/usr/bin/bash
# export GOPATH=C:/Users/ka-ha/go
# export GOMODCACHE=C:/Users/ka-ha/go/pkg/mod
# Go パラメータ
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMODTIDY=$(GOCMD) mod tidy
BINARY_NAME=hugo-docs-i18n.exe
CONFIG_NAME=hugo-docs-i18n.yaml
# for Test data
TEST_DIR=content/ja
DIST_DIR=dist
REPO_VER=version.json
CUR_VER=doci18n/version.json
TEST_VER=doci18n/testdata/version.json
TEST_EN=locale/testdata/l10n.en.json
CUR_EN=locale/l10n.en.json
REPO_EN=data/i18n/l10n.en.json
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
		cp -f $(REPO_VER) $(CUR_VER)
		cp -f $(REPO_EN) $(CUR_EN)
		$(GOBUILD) -o $(BINARY_NAME) -v -buildvcs=false
test: test-doci18n test-locale test-cmd
#		$(GOTEST) -v ./...
# to execute test-doci18n need to copy doci18n/testdata/version.go to doci18n/testdata/ 
test-doci18n:
		cp $(TEST_VER) $(CUR_VER)
		$(GOTEST) -v ./doci18n
test-locale:
		cp $(TEST_EN) $(CUR_EN)
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
		@echo 'if "help" tag exists, "git tag -d help"'
		@echo 'execute "git push origin --tags"'
#		git push origin --tags
git-remote-tag:
		git ls-remote https://github.com/juggernautjp/hugo-docs-i18n/ | tail -1
git-fix:
		git config --global --add safe.directory F:/Docs.repo/hugo-docs-i18n
now: $(BINARY_NAME)
		@echo $(NOW) $(GIT_VER) "$(VER_MSG)"
install:
		go install github.com/juggernautjp/hugo-docs-i18n@v0.3.3
