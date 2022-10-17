/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Markdown file FrontMatter

Get data from Front Matter of Hugo Markdown file with:
	- [hugo/commands/convert.go]()
*/

package doci18n

import (
	"io"
	"os"
	"bytes"
	"fmt"
	// "log"
	"strings"
	// "unsafe"
	// "encoding/json"
	"time"

	// "gopkg.in/yaml.v2"
	// "github.com/BurntSushi/toml"
	// toml "github.com/pelletier/go-toml/v2"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/spf13/cast"
)

// File is exist?
func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// File is empty?
func IsEmpty(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.Size() == 0
}

// Is directory?
func IsDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

// Is regular file?
func IsRegularFile(filename string) bool {
	if filename == "" || !IsExist(filename) || IsEmpty(filename) || IsDir(filename) {
		return false
	}
	return true
}

// Read content file, return FrontMatter and Content
func ReadContentFile(infn string) (pageparser.ContentFrontMatter, error) {
	var pf pageparser.ContentFrontMatter
	var err error
	var contentBytes []byte

	// Check if the file is file and not zero-size
	if !IsRegularFile(infn) {
		return pf, fmt.Errorf("File is not regular: %s", infn)
	}

	// Check if the specified file is content file
	if !files.IsContentFile(infn) {
		return pf, fmt.Errorf("File is not known content format: %s", infn)
	}

	// ReadFile 
	contentBytes, err = os.ReadFile(infn)
	if err != nil {
		// return pf, fmt.Errorf("Failed to read file %q: %w", infn, err)
		return pf, err
	}

	// var pf ContentFrontMatter (defined in parser/pageparser/pageparser.go)
	pf, err = pageparser.ParseFrontMatterAndContent(bytes.NewReader(contentBytes))
	if err != nil {
		// return pf, fmt.Errorf("Failed to parse file %q: %w", infn, err)
		return pf, err
	}

	// better handling of dates in formats that don't have support for them
	if pf.FrontMatterFormat == metadecoders.JSON || pf.FrontMatterFormat == metadecoders.YAML || pf.FrontMatterFormat == metadecoders.TOML {
		for k, v := range pf.FrontMatter {
			switch vv := v.(type) {
			case time.Time:
				pf.FrontMatter[k] = vv.Format(time.RFC3339)
			}
		}
	} else {
		// pageparser.ParseFrontMatterAndContent() return pf.Content == nil, if the file has no FrontMatter.
		// So, ReadContentFile() set pf.Content with contentBytes that returned by os.ReadFile().
		pf.Content = contentBytes
	}
	return pf, nil
}

// return if the file is draft or not.
func IsDraftFile(infn string) (bool, string, error) {
	// ReadContentFile()
	pf, err := ReadContentFile(infn)
	if err != nil {
		return false, "", err
	}

	// skip if the page is draft (i.e. "draft" of FrontMatter == true)
	// copy the page with setting field "draft" of FrontMatter = true
	var isDraft bool = false
	var sTitle string = ""
	for k, v := range pf.FrontMatter {
		loki := strings.ToLower(k)
		switch loki {
		case "draft":
			isDraft = cast.ToBool(v)
		case "title":
			sTitle = cast.ToString(v)
		}
	}
	// if this page is draft, return true (the page is draft).
	return isDraft, sTitle, nil
}

// copy not-draft file from <infn> to <outfn>,
// with setting which draft = true for translation. 
func CopyContentFile(infn, outfn string) (bool, error) {
	// if outfn is "", return false/error.
	if outfn == "" {
		return false, fmt.Errorf("Dst file is not specified")
	}
	// if file is already exist
	/*
	if IsExist(outfn) {
		return false, fmt.Errorf("File already exist")
	}
	*/

	// ReadContentFile()
	pf, err := ReadContentFile(infn)
	if err != nil {
		return false, err
	}

	// skip if the page is draft (i.e. "draft" of FrontMatter == true)
	// copy the page with setting field "draft" of FrontMatter = true
	var isDraft bool = false
	var i int = 0
	for k, v := range pf.FrontMatter {
		i++
		loki := strings.ToLower(k)
		switch loki {
		case "draft":
			isDraft = cast.ToBool(v)
			// set "draft" value of FrontMatter to true
			pf.FrontMatter[k] = true
		}
	}
	// if this page is draft, return true (the page is draft).
	if isDraft {
		return false, nil
	}
	// if pf.FrontMatter dose not have "draft" field, add "draft: true"
	if i == len(pf.FrontMatter) {
		pf.FrontMatter["draft"] = true
	}

	// Marshal FrontMatter
	var newContent bytes.Buffer
	if pf.FrontMatterFormat == "" {
		// if the file has no FrontMatter, write YAML FrontMatter with "draft: true"
		pf.FrontMatterFormat = metadecoders.YAML
	}
	err = parser.InterfaceToFrontMatter(pf.FrontMatter, pf.FrontMatterFormat, &newContent)
	if err != nil {
		return false, err
	}
	newContent.Write(pf.Content)

	// output FrontMatter and Content to outfn
	outfile, err := os.Create(outfn)
	if err != nil {
		// return false, fmt.Errorf("Failed to create file %q: %w", outfn, err)
		return false, err
	}
	// close the file
	defer outfile.Close()
	// write frontmatter and content to file
	fmt.Fprint(outfile, &newContent)

	return true, nil
}

// Copy not-content file, including CSS, jpeg, etc. 
func CopyNotContentFile(infn, outfn string) (bool, error) {
	// Check if the file is file and not zero-size
	if !IsRegularFile(infn) {
		return false, fmt.Errorf("File is not regular")
	}

	// Check if the specified file is content file
	if files.IsContentFile(infn) {
		return false, nil
	}

	// if the file is already exist, overwrite
	/*  
	if IsExist(outfn) {
		return false
	}
	*/

	// copy file
	src, err := os.Open(infn)
	if err != nil {
		return false, err
	}
	defer src.Close()

	dst, err := os.Create(outfn)
	if err != nil {
		return false, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return false, err
	}
	return true, nil
}
