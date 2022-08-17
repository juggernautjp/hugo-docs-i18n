/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n Frontmatter Test
*/

package doci18n

import (
	"embed"
	// "fmt"
	// "io/fs"
	// "os"
	// "log"
	"encoding/json"
	"path/filepath"
	// "reflect"
	"testing"

)

// Sample Front Matter
type Sample struct {
  Title string `json:"title"`
  LinkTitle string `json:"linktitle"`
  Description string `json:"description"`
  Date string `json:"date"`
  Categories []string `json:"categories"`
  Keywords []string `json:"keywords"`
  // "menu": { "docs": { "parent": "modules", "weight": 50 } }
  Weight int `json:"weight"`
  SectionWeight int `json:"sections_weight"`
  Draft bool `json:"draft"`
	LastMod   string `json:"lastmod"`
  Aliases []string `json:"aliases"`
  Toc bool `json:"toc"`
}

//go:embed testdata/en/NO_FrontMatter.md
var static embed.FS

//go:embed testdata/en/Sample.json
var sampleBytes []byte

// Test for ReadContentFile()
func TestReadContentFile(t *testing.T) {
	var matter hugoFrontMatter
	// Test data #1: Front Matter from JSON
	var wantFM Sample
	if err := json.Unmarshal(sampleBytes, &wantFM); err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", wantFM)

	// Test data #2: Contents from Markdown file
	b, err := static.ReadFile("testdata/en/NO_FrontMatter.md")
	if err != nil {
		panic(err)
	}
	wantRest := string(b)
	// fmt.Printf("%s\n", wantRest)

	// Test data #3: File name to read for test
	testCases := []string {
		"JSON_FrontMatter.md",
		"TOML_FrontMatter.md",
		"YAML_FrontMatter.md",
	}
	// dir, _ := os.Getwd()
	const dir = "testdata/en"

	// Test each Field of Struct
	failFunc := func(in string, want, act interface{}) {
		t.Fatalf("Key: `%s`\n\nNot equal: \n"+
			"expected: %v\n"+
			"actual  : %v", in, want, act)
	}

	checkFunc := func(wantFM *Sample, gotFM *hugoFrontMatter) {
		if wantFM.Title != gotFM.Title {
			failFunc("Title: ", wantFM.Title, gotFM.Title)
		}
		if wantFM.Date != gotFM.Date {
			failFunc("Date: ", wantFM.Date, gotFM.Date)
		}
		if wantFM.LastMod != gotFM.LastMod {
			failFunc("LastMod: ", wantFM.LastMod, gotFM.LastMod)
		}
		if wantFM.Draft != gotFM.Draft {
			failFunc("Draft: ", wantFM.Draft, gotFM.Draft)
		}
	}

	// Run each test data
	for _, tc := range testCases {
		fname := filepath.Join(dir, tc)
		// Test Run
		t.Run(tc, func(t *testing.T) {
			gotRest := ReadContentFile(fname, &matter)
			checkFunc(&wantFM, &matter)
			if gotRest != wantRest {
				t.Errorf("Rest Contents: \n\ngot : %v\nwant: %v\n", gotRest, wantRest)
			}
		})
	}
	/*
	want := []string{
		filepath.Join(dir, "dir_b", "file_b"),
		filepath.Join(dir, "dir_c", "file_c"),
		filepath.Join(dir, "file_a"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadContentFile():\ngot  %v\nwant %v", got, want)
	}
	*/
}