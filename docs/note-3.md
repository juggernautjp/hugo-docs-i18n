
# note-3

[TOC]

## listlangpage.go

翻訳中のファイルを探索し、翻訳済みと未翻訳のファイルを判定する関数の実装。

github.com/gohugoio/hugo/commands/list.go を参考にした実装方法は、
github.com/gohugoio/hugo/commands 以下の複数の Go 言語ファイルを必要とするため、
実装ができないことが判明した。

このため、ファイルから FrontMatter を解析し、draft フィールドの値を取り出す関数を実装する。


## reflection

```go
		for k, v := range fm {
			loki := strings.ToLower(k)
			switch loki {
			case "title":
				gotFM.Title = cast.ToString(v)
			case "date":
				gotFM.Date = cast.ToString(v)
				/*
				pv := reflect.ValueOf(v)
				if pv.Kind() == reflect.String {
					pv.Set(gotFM.Date)
				}
				*/
			case "weight":
				gotFM.Weight = cast.ToInt(v)
			case "draft":
				gotFM.Draft = cast.ToBool(v)
			}
		}
```