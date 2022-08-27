
# hugo-docs-ja

[TOC]

## コマンドの実行方法

```ps1
# Go Module の初期化
$ go mod init hugo-docs-i18n

# Cobra CLI アプリケーションのインストール
$ go install github.com/spf13/cobra-cli@latest

# Cobra CLI アプリケーションの起動 (Cobra の初期化)
$ cobra-cli init 
$ go run main.go
$ go run main.go --help

# サブコマンドの追加
$ cobra-cli add version
$ cobra-cli add init
$ cobra-cli add update
$ cobra-cli add debug
$ cobra-cli add output -p 'debugCmd'
$ cobra-cli add reset -p 'debugCmd'

# コンパイル
$ go build main.go
$ mv main.exe hugo-docs-i18n.exe
```


## go test コマンドの実行方法

```ps1
# locale ライブラリ (locale ディレクトリ以下のファイル) のテスト実行 
$ go test -v ./locale
# すべてのファイルの実行
$ go test -v ./...
```



## Go モジュールのインストール

```ps1
$ go get github.com/gohugoio/hugo/hugofs/files
$ go get github.com/gohugoio/hugo/parser/pageparser
```

Go モジュール

```ps1
$ go mod tidy
$ go mod verify
```



## debug コマンド

- *`debug reset` コマンド実行前*

```go
$ go run main.go debug output
output called
Aliases:
map[string]string{}
Override:
map[string]interface {}{}
PFlags:
map[string]viper.FlagValue{"code":viper.pflagValue{flag:(*pflag.Flag)(0xc00016c280)}, "content-dir":viper.pflagValue{flag:(*pflag.Flag)(0xc000001d60)}, "lang-name":viper.pflagValue{flag:(*pflag.Flag)(0xc000001e00)}, "language":viper.pflagValue{flag:(*pflag.Flag)(0xc00016c1e0)}, "locale":viper.pflagValue{flag:(*pflag.Flag)(0xc000001cc0)}, "locale-db":viper.pflagValue{flag:(*pflag.Flag)(0xc000001b80)}, "src-md":viper.pflagValue{flag:(*pflag.Flag)(0xc000001c20)}, "time-format-blog":viper.pflagValue{flag:(*pflag.Flag)(0xc00016c000)}, "time-format-default":viper.pflagValue{flag:(*pflag.Flag)(0xc000001f40)}, "weight":viper.pflagValue{flag:(*pflag.Flag)(0xc000001ea0)}}
Env:
map[string][]string{}
Key/Value Store:
map[string]interface {}{}
Config:
map[string]interface {}{}
Defaults:
map[string]interface {}{}
```

- *`debug reset` コマンド実行後*

```go
$ go run main.go debug output
output called
Aliases:
map[string]string{}
Override:
map[string]interface {}{}
PFlags:
map[string]viper.FlagValue{"code":viper.pflagValue{flag:(*pflag.Flag)(0xc00016c1e0)}, "content-dir":viper.pflagValue{flag:(*pflag.Flag)(0xc000001cc0)}, "lang-name":viper.pflagValue{flag:(*pflag.Flag)(0xc000001d60)}, "language":viper.pflagValue{flag:(*pflag.Flag)(0xc00016c140)}, "locale":viper.pflagValue{flag:(*pflag.Flag)(0xc000001c20)}, "locale-db":viper.pflagValue{flag:(*pflag.Flag)(0xc000001ae0)}, "src-md":viper.pflagValue{flag:(*pflag.Flag)(0xc000001b80)}, "time-format-blog":viper.pflagValue{flag:(*pflag.Flag)(0xc000001f40)}, "time-format-default":viper.pflagValue{flag:(*pflag.Flag)(0xc000001ea0)}, "weight":viper.pflagValue{flag:(*pflag.Flag)(0xc000001e00)}}
Env:
map[string][]string{}
Key/Value Store:
map[string]interface {}{}
Config:
map[string]interface {}{}
Defaults:
map[string]interface {}{}
```
