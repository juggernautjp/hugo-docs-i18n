/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n searchdb Test
*/

package locale

import (
	// "embed"
	// "fmt"
	// "io/fs"
	// "os"
	// "log"
	///"encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)

// Test for SearchLocaleFile()
func TestSearchLocaleFile(t *testing.T) {
	const dir = "testdata"
	const infn = "ISO_639-1.json"
	infname := filepath.Join(dir, infn)
	// test case for lang/code pair
	testCasesLang := []string {
		"Ar",
		"Hong",
	}
	var wantLangPair = [...][]LangPair {
		{
			{ Lang: "Arabic (Algeria)", Code: "ar-dz" },
			{ Lang: "Arabic (Bahrain)", Code: "ar-bh" },
			{ Lang: "Arabic (U.A.E.)", Code: "ar-ae" },
		},
		{
			{ Lang: "Chinese (Hong Kong)", Code: "zh-hk" },
		},
	}
	testCasesCode := []string {
		"ar",
		"md",
	}
	var wantLangPair2 = [...][]LangPair {
		{
			{ Lang: "Arabic (Algeria)", Code: "ar-dz" },
			{ Lang: "Arabic (Bahrain)", Code: "ar-bh" },
			{ Lang: "Arabic (U.A.E.)", Code: "ar-ae" },
		},
		{
			{ Lang: "Romanian (Republic of Moldova)", Code: "ro-md" },
			{ Lang: "Russian (Republic of Moldova)", Code: "ru-md" },
		},
	}

	// Case 1: search by language
	t.Run("Case 1: by language", func(t *testing.T) {
		for k, lg := range testCasesLang {
			// Test Run
			got, err := SearchLocaleFile(infname, lg, "")
			if err != nil {
				t.Fatalf("Failed to search %s: %v", lg, err)
			}
			if !reflect.DeepEqual(got, wantLangPair[k]) {
				t.Errorf("SearchLocaleFile(%s):\ngot  %v\nwant %v", lg, got, wantLangPair[k])
			}
		}
	})

	// Case 2: search by code
	t.Run("Case 2: by code", func(t *testing.T) {
		for k, cd := range testCasesCode {
			// Test Run
			got, err := SearchLocaleFile(infname, "", cd)
			if err != nil {
				t.Fatalf("Failed to search %s: %v", cd, err)
			}
			if !reflect.DeepEqual(got, wantLangPair2[k]) {
				t.Errorf("SearchLocaleFile(%s):\ngot  %v\nwant %v", cd, got, wantLangPair2[k])
			}
		}
	})
}
