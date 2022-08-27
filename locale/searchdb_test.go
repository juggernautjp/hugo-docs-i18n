/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n searchdb Test
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
		"ar",
		"Hong",
	}
	wantLangPair[0] = []LangPair {
		{ Lang: "Arabic (Algeria)", Code: "ar-dz" },
		{ Lang: "Arabic (Bahrain)", Code: "ar-bh" },
		{ Lang: "Arabic (U.A.E.)", Code: "ar-ae" },
	}
	wantLangPair[1] = []LangPair {
		{ Lang: "Chinese (Hong Kong)", Code: "zh-hk" },
	}

	// Case 1: search by language
	t.Run("Case 1: by language", func(t *testing.T) {
		for k, lg := range testCasesLang {
			// Test Run
			got, err := SearchLocaleFile(infnam, lg, "")
			if err != nil {
				t.Fatalf("Failed to search %s: %v", lg, err)
			}
			if !reflect.DeepEqual(got, wantLangPair[k]) {
				t.Errorf("WalkDir3():\ngot  %v\nwant %v", got, wantLangPair[0])
			}
		}
	})
}
