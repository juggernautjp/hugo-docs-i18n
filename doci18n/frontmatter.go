/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Markdown file FrontMatter

Get data from Front Matter of Hugo Markdown file with:
	- [hugo/commands/convert.go]()
*/

package doci18n

import (
	"os"
	"bytes"
	// "fmt"
	"log"
	"strings"
	"unsafe"
	"encoding/json"
	"time"

	// "gopkg.in/yaml.v2"
	// "github.com/BurntSushi/toml"
	// toml "github.com/pelletier/go-toml/v2"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/spf13/cast"
)

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

// get Draft from Front Matter
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
}

// Read Front Matter and Contents from the file `fname`
func ReadContentFile(fname string) (ContentFrontMatter, error) {
	// 
	if !files.IsContentFile(fname) {
		return nil, fmt.Errorf("target path %q is not a known content format", fname)
	}

	// ReadFile 
	contentBytes, err := os.ReadFile(fname)
	if err != nil {
		log.Fatalf("Read file error:", fname)
		return nil, err
	}

	// var pf ContentFrontMatter (defined in parser/pageparser/pageparser.go)
	pf, err := pageparser.ParseFrontMatterAndContent(bytes.NewReader(contentBytes))
	if err != nil {
		log.Fatalf("Parse file error:", fname)
		return nil, err
	}

	// better handling of dates in formats that don't have support for them
	if pf.FrontMatterFormat == metadecoders.JSON || pf.FrontMatterFormat == metadecoders.YAML || pf.FrontMatterFormat == metadecoders.TOML {
		for k, v := range pf.FrontMatter {
			switch vv := v.(type) {
			case time.Time:
				pf.FrontMatter[k] = vv.Format(time.RFC3339)
			}
		}
	}

	// convert []byte to string
	// rest = *(*string)(unsafe.Pointer(&ret))
	return pf, nil

	// Output:
	// {Name:frontmatter Tags:[go yaml json toml]}
	// rest of the content
}
