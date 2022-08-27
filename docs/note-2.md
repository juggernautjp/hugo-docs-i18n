
# note-2

## FrontMatter

```go
// github.com/gohugoio/hugo/hugolib/page_meta.go
type pageMeta struct {
  draft     bool
  title     string
}

func (pm *pageMeta) Draft() bool {
  return pm.draft
}

func (pm *pageMeta) Title() string {
  return pm.title
}

// get each field data from FrontMatter
func (pm *pageMeta) setMetadata(frontmatter map[string]any) error {
  // var draft *bool
  for k, v := range frontmatter {
    loki := strings.ToLower(k)
    switch loki {
    case "title":
      pm.title = cast.ToString(v)
    case "draft":
      // draft = new(bool)
      // *draft = cast.ToBool(v)
      pm.draft = cast.ToBool(v)
    }
  }
  return nil
}
```



## convert `[]byte` to string

```go
// convert []byte to string
rest = *(*string)(unsafe.Pointer(&ret))
```


## Output comment

```go
func name () {
  // Output:
  // {Name:frontmatter Tags:[go yaml json toml]}
  // rest of the content
}
```