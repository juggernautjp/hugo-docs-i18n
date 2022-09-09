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
	"log"
	"encoding/json"
	"path/filepath"
	// "reflect"
	"testing"

	"hugo-docs-i18n/doci18n"
)

// go:embed testdata/ISO_639-1.json
// var sampleBytes []byte

// Locale (lang/code pair) data
type langPair struct {
	lang string `json:"lang"`
	code string `json:"code"`
}

type langJSON struct {
	locale []langPair `json:"locale"`
}

// Test for ConvertLocaleFile()
func TestConvertLocaleFile(t *testing.T) {
	const dir = "testdata"
	const infn = "ISO_639-1.json"
	infname := filepath.Join(dir, infn)
	var sampleBytes []byte
	sampleBytes, err := os.ReadFile(infname)
	if err != nil {
		log.Fatal(err)
	}
	// Test data #1: Locale (lang/code pair) data from JSON
	var wantJson langJSON
	if err := json.Unmarshal(sampleBytes, &wantJson); err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v\n", wantJson)

	// Prepare for Test data #2
	const outfn = "out.json"
	const inmdfn = "ISO_639-1.md"
	inmdfname := filepath.Join(dir, inmdfn) 
	outfname := filepath.Join(dir, outfn)
	if doci18n.IsExist(outfname) {
		os.Remove(outfname)
	}

	// Test data #2: Locale (lang/code pair) data from Markdown file
	if err := ConvertLocaleFile(inmdfname, outfname); err != nil {
		log.Fatal(err)
	}
	
	// Test data #3: Locale (lang/code pair) data from JSON output by #2
	contentBytes, err := os.ReadFile(outfname)
	if err != nil {
		log.Fatal(err)
	}
	var gotJson langJSON
	if err := json.Unmarshal(contentBytes, &gotJson); err != nil {
		log.Fatal(err)
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
		t.Run("locale lang/code pairs", func(t *testing.T) {
			checkFunc(&lp, &gotJson.locale[i])
		})
	}
}