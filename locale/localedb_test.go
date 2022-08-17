/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n locale Test
*/

package locale

import (
	// "embed"
	// "fmt"
	// "io/fs"
	"os"
	// "log"
	"encoding/json"
	"path/filepath"
	// "reflect"
	"testing"
)

// go:embed testdata/ISO_639-1.json
// var sampleBytes []byte

// Sample Front Matter
type langPair struct {
  lang string `json:"lang"`
  code string `json:"code"`
}

type Sample struct {
  locale []langPair `json:"locale"`
}

// Test for ReadLocaleFile()
func TestReadLocaleFile(t *testing.T) {
	const dir = "testdata"
	const infn = "ISO_639-1.json"
	infname := filepath.Join(dir, infn)
	var sampleBytes []byte
	sampleBytes, err := os.ReadFile(infname)
	if err != nil {
		panic(err)
	}
	// Test data #1: Locale (lang/code pair) date from JSON
	var wantJson Sample
	if err := json.Unmarshal(sampleBytes, &wantJson); err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", wantJson)

	// Test data #2: Locale (lang/code pair) date from Markdown file
	const outfn = "out.json"
	outfname := filepath.Join(dir, outfn)
	if err := ReadLocaleFile("testdata/ISO_639-1.md", outfname); err != nil {
		panic(err)
	}
	
	// Test data #3: Locale (lang/code pair) date from JSON output by #2
	contentBytes, err := os.ReadFile(outfname)
	if err != nil {
		panic(err)
	}
	var gotJson Sample
	if err := json.Unmarshal(contentBytes, &gotJson); err != nil {
		panic(err)
	}

	// Verify wantJson == gotJson
	failFunc := func(in string, want, act interface{}) {
		t.Fatalf("Key: `%s`\n\nNot equal: \n"+
			"expected: %v\n"+
			"actual  : %v", in, want, act)
	}

	checkFunc := func(wantJson, gotJson *langPair) {
		if wantJson.lang != gotJson.lang {
			failFunc("lang: ", wantJson.lang, gotJson.lang)
		}
		if wantJson.code != gotJson.code {
			failFunc("code: ", wantJson.code, gotJson.code)
		}
	}

	// Run each test data
	for i, lp := range wantJson.locale {
		// Test Run
		t.Run("locale rest", func(t *testing.T) {
			checkFunc(&lp, &gotJson.locale[i])
		})
	}
}