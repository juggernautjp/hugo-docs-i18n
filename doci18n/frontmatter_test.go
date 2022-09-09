/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n Frontmatter Test
*/

package doci18n

import (
	"embed"
	// "fmt"
	// "io/fs"
	"os"
	"log"
	"strings"
	// "time"
	"encoding/json"
	"path/filepath"
	// "reflect"
	"testing"

	"github.com/spf13/cast"
)

// Sample FrontMatter
type Sample struct {
	Title string `json:"title"`
	// LinkTitle string `json:"linktitle"`
	// Description string `json:"description"`
	Date string `json:"date"`
	// Categories []string `json:"categories"`
	// Keywords []string `json:"keywords"`
	// "menu": { "docs": { "parent": "modules", "weight": 50 } }
	Weight int `json:"weight"`
	// SectionWeight int `json:"sections_weight"`
	Draft bool `json:"draft"`
	// LastMod   string `json:"lastmod"`
	// Aliases []string `json:"aliases"`
	// Toc bool `json:"toc"`
}

// Hugo FrontMatter definition to necessary for i18n
type frontMatter struct {
	Title     string     `yaml:"title" toml:"title" json:"title"`
	Date      string     `yaml:"date" toml:"date" json:"date"`
	Weight    int        `yaml:"weight" toml:"weight" json:"weight"`
	// tDate     time.Time  `toml:"date"`
	Draft     bool       `yaml:"draft" toml:"draft" json:"draft"`
	// Tags      []string `yaml:"tags" toml:"tags" json:"tags"`
	// Metadata metadata  `yaml:"metadata" toml:"metadata" json:"metadata"`
}

//go:embed testdata/en/NO_FrontMatter.md
var static embed.FS

//go:embed testdata/en/Sample.json
var sampleBytes []byte

// createEmptyFile() create a zero-size file
func createEmptyFile(file string) {
	if err := os.WriteFile(file, []byte(""), 0666); err != nil {
		log.Fatal(err)
	}
}
// createNotEmptyFile() create a zero-size file
func createNotEmptyFile(file string) {
	if err := os.WriteFile(file, []byte("content"), 0666); err != nil {
		log.Fatal(err)
	}
}

// Test for IsExist()
func TestIsExist(t *testing.T) {
	dir := t.TempDir()
	testCases := map[string]bool {
		"file_a.md": true,
		"file_b.md": false,
	}
	t.Run("IsExist()", func(t *testing.T) {
		for fn, want := range testCases {
			infname := filepath.Join(dir, fn)
			if want {
				createEmptyFile(infname)
			}
			got := IsExist(infname)
			if got != want {
				t.Errorf("File is exist?: %s\n\ngot : %v\nwant: %v\n", infname, got, want)
			}
		}
	})
}

// Test for IsEmpty()
func TestIsEmpty(t *testing.T) {
	dir := t.TempDir()
	testCases := map[string]bool {
		"file_a.md": true,
		"file_b.md": false,
		"file_c.md": false,
	}
	t.Run("IsEmpty()", func(t *testing.T) {
		for fn, want := range testCases {
			infname := filepath.Join(dir, fn)
			if want {
				createEmptyFile(infname)
			}
			if fn == "file_c.md" {
				createNotEmptyFile(infname)
			}
			got := IsEmpty(infname)
			if got != want {
				t.Errorf("File is empty?: %s\n\ngot : %v\nwant: %v\n", infname, got, want)
			}
		}
	})
}

// Test for IsDir()
func TestIsDir(t *testing.T) {
	dir := t.TempDir()
	testCases := map[string]bool {
		"dir_a": true,
		"dir_b": false,
		"file_b.md": false,
	}
	t.Run("IsDir()", func(t *testing.T) {
		for fn, want := range testCases {
			infname := filepath.Join(dir, fn)
			if want {
				os.MkdirAll(infname, 0755)
			}
			if filepath.Ext(fn) != "" {
				createNotEmptyFile(infname)
			}
			got := IsDir(infname)
			if got != want {
				t.Errorf("Is directory?: %s\n\ngot : %v\nwant: %v\n", infname, got, want)
			}
		}
	})
}

// Test for ReadContentFile()
func TestReadContentFile(t *testing.T) {
	// Test data #1: FrontMatter from JSON
	var wantFM Sample
	if err := json.Unmarshal(sampleBytes, &wantFM); err != nil {
		t.Fatalf("Failed to read FrontMatter JSON: %v", err)
	}
	// fmt.Printf("%+v\n", wantFM)

	// Test data #2: Content from Markdown file
	const dir = "testdata/en"
	infn := filepath.Join(dir, "NO_FrontMatter.md")
	// b, err := static.ReadFile(infn)
	b, err := os.ReadFile(infn)
	if err != nil {
		t.Fatalf("Failed to read Content Markdown: %v", err)
	}
	wantRest := string(b)
	// fmt.Printf("%s\n", wantRest)

	// Test data #4: File name to read for copy file
	testCases := []string {
		"JSON_FrontMatter.md",
		"TOML_FrontMatter.md",
		"YAML_FrontMatter.md",
	}

	// Test each Field of Struct
	failFunc := func(in string, want, act interface{}, fn string) {
		t.Errorf("Key: `%s` in %s\n\nNot equal: \n"+
			"expected: %v\n"+
			"actual  : %v\n", in, fn, want, act)
	}

	checkFunc := func(wantFM *Sample, fm map[string]any, fn string) {
		var gotFM frontMatter
		for k, v := range fm {
			loki := strings.ToLower(k)
			switch loki {
			case "title":
				gotFM.Title = cast.ToString(v)
			case "date":
				gotFM.Date = cast.ToString(v)
			case "weight":
				gotFM.Weight = cast.ToInt(v)
			case "draft":
				gotFM.Draft = cast.ToBool(v)
			}
		}

		if wantFM.Title != gotFM.Title {
			failFunc("Title", wantFM.Title, gotFM.Title, fn)
		}
		if wantFM.Date != gotFM.Date {
			failFunc("Date", wantFM.Date, gotFM.Date, fn)
		}
		if wantFM.Weight != gotFM.Weight {
			failFunc("Weight", wantFM.Weight , gotFM.Weight, fn)
		}
		if wantFM.Draft != gotFM.Draft {
			failFunc("Draft", wantFM.Draft, gotFM.Draft, fn)
		}
	}

	// Case 1: 
	t.Run("Case 1: TestReadContentFile()", func(t *testing.T) {
		for _, fn := range testCases {
			infname := filepath.Join(dir, fn)
			// Test Run
			cfm, err := ReadContentFile(infname)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", infn, err)
			}
			checkFunc(&wantFM, cfm.FrontMatter, fn)
			gotRest := string(cfm.Content)
			if gotRest != wantRest {
				t.Errorf("Rest Contents: %s\n\ngot : %v\nwant: %v\n", fn, gotRest, wantRest)
			}
		}
	})
}

// Test for IsDraftFile()
func TestIsDraftFile(t *testing.T) {
	// Test data #3: File name to read for draft test
	const dir = "testdata"
	testCases := []struct {
		path string
		title string
		draft bool
	}{
		{
			path: "YAML_Draft.md",
			title: "YAML Draft Markdown",
			draft: true,
		},
		{
			path: "YAML_NotDraft.md",
			title: "YAML NotDraft Markdown",
			draft: false,
		},
	}

	// Case 1: 
	t.Run("Case 2: TestIsDraftFile()", func(t *testing.T) {
		for _, want := range testCases {
			infname := filepath.Join(dir, want.path)
			// Test Run
			gotDraft, gotTitle, err := IsDraftFile(infname)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", infname, err)
			}
			if gotDraft != want.draft {
				t.Errorf("File is draft?: %s\n\ngot : %v\nwant: %v\n", infname, gotDraft, want.draft)
			}
			if gotTitle != want.title {
				t.Errorf("Title: %s\n\ngot : %v\nwant: %v\n", infname, gotTitle, want.title)
			}
		}
	})
}

// Test for CopyContentFile()
func TestCopyContentFile(t *testing.T) {
	// Test data #3: File name to read for draft test 
	const dir = "testdata"
	testCases := map[string]bool {
		"YAML_Draft.md": false,
		"YAML_NotDraft.md": true,
	}
	outdir := t.TempDir()
	/*
	const outdir = ""
	if !IsExist(outdir) {
		os.Mkdir(outdir, 0755)
	}
	*/

	// Case 3: 
	t.Run("Case 3: TestCopyContentFile()", func(t *testing.T) {
		for fn, want := range testCases {
			infname := filepath.Join(dir, fn)
			outfname := filepath.Join(outdir, fn)
			// if outfname is exist, remove it
			if IsExist(outfname) {
				os.Remove(outfname)
			}

			// Test Run
			got, err := CopyContentFile(infname, outfname)
			if err != nil {
				t.Fatalf("Failed to copy file %s: %v", infname, err)
			}
			if got != want {
				t.Errorf("Faild to copy %s to %s\n", infname, outfname)
			}
			if got {
				// File copied
				if !IsExist(outfname) {
					t.Errorf("Output file should be copied: %s\n", outfname)
				} else {
					got, _, err := IsDraftFile(outfname)
					if err != nil {
						t.Fatalf("Failed to read copied file %s: %v", outfname, err)
					}
					if !got {
						t.Errorf("Output file is not draft: %s\n", infname)
					}
				}
			} else {
				// File not copied
				if IsExist(outfname) {
					t.Errorf("Output file should not be copied: %s\n", outfname)
				}
			}
		}
	})
}

// Test for CopyNotContentFile()
func TestCopyNotContentFile(t *testing.T) {
	// Test data #3: File name to read for draft test 
	const dir = "testdata/zn"
	testCases := map[string]bool {
		"images/hugo-with-nanobox.png": true,
		"css/style.css": true,
		"dir_a/draft_a.md": false,
		"dir_a/nofile.md": false,
	}
	outdir := t.TempDir()
	// Case 3: 
	t.Run("Case 4: TestCopyNotContentFile()", func(t *testing.T) {
		for p, want := range testCases {
			_, fn := filepath.Split(p)
			infname := filepath.Join(dir, p)
			outfname := filepath.Join(outdir, fn)
			if IsExist(outfname) {
				os.Remove(outfname)
			}
			// Test Run
			got := CopyNotContentFile(infname, outfname)
			if got != want {
				t.Errorf("Faild to copy %s to %s\n", infname, outfname)
			}
		}
	})
}
