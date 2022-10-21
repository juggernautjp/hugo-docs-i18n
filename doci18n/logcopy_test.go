/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n logcopy Test
*/

package doci18n

import (
	// "fmt"
	// "io/fs"
	"os"
	"log"
	// "strings"
	"path/filepath"
	"encoding/json"
	"reflect"
	"testing"
)


// Test for InitLogJSON()
func TestInitLogJSON(t *testing.T) {
	cd := GetCurrentString()
	sd := "testdata/zn"
	dd := "testdata/ko"
	wantSJ := SummaryJSON{ Date: cd, SrcDir: sd, DstDir: dd,}
	wantFJ := FileJSON{ ContentFiles: []CopyFile{}, NotContentFiles: []CopyFile{}, }
	t.Run("NewLogJSON", func(t *testing.T) {
		InitLogJSON(cd, sd, dd)
		gotLJ := GetLogJSON()
		gotSJ := gotLJ.Header
		gotFJ := gotLJ.Paths
		if !reflect.DeepEqual(gotSJ, wantSJ) {
			t.Errorf("Compare SummaryJSON:\ngot %v\nwant %v", gotSJ, wantSJ)
		}
		if !reflect.DeepEqual(gotFJ, wantFJ) {
			t.Errorf("Compare FileJSON:\ngot %v\nwant %v", gotFJ, wantFJ)
		}
	})
}

// Test for AddContentFile()
func TestAddContentFile(t *testing.T) {
	cd := GetCurrentString()
	sd := "testdata/zn"
	dir := "testdata/ko"
	wantFJ := FileJSON{
		ContentFiles: []CopyFile{
			{ Path: filepath.Join(dir, "dir_a/draft_a.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "dir_a/notdraft_a.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "dir_b/draft_b.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "dir_b/notdraft_b.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "draft_c.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "myshowcase/no_draft.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "myshowcase/no_frontmatter.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "notdraft_c.md"), Result: NotUpdated, },
		},
		NotContentFiles: []CopyFile{
			{ Path: filepath.Join(dir, "css/style.css"), Result: Copied, },
			{ Path: filepath.Join(dir, "myshowcase/hugo-with-nanobox.png"), Result: Overwrited, },
		},
	}

	// 
	InitLogJSON(cd, sd, dir)
	for _, want := range wantFJ.ContentFiles {
		AddContentFile(want.Path, want.Result)
	}
	for _, want := range wantFJ.NotContentFiles {
		AddNotContentFile(want.Path, want.Result)
	}

	// compare JSON
	t.Run("AddContentFile", func(t *testing.T) {
		gotLJ := GetLogJSON()
		gotFJ := gotLJ.Paths
		if !reflect.DeepEqual(gotFJ, wantFJ) {
			t.Errorf("Compare JSON:\ngot %v\nwant %v", gotFJ, wantFJ)
		}
	})
}

// Test for SaveLogJSON()
func TestSaveLogJSON(t *testing.T) {
	cd := GetCurrentString()
	sd := "testdata/zn"
	dir := "testdata/ko"
	wantFJ := FileJSON{
		ContentFiles: []CopyFile{
			{ Path: filepath.Join(dir, "dir_a/draft_a.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "dir_a/notdraft_a.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "dir_b/draft_b.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "dir_b/notdraft_b.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "draft_c.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "myshowcase/no_draft.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "myshowcase/no_frontmatter.md"), Result: Copied, },
			{ Path: filepath.Join(dir, "notdraft_c.md"), Result: NotUpdated, },
		},
		NotContentFiles: []CopyFile{
			{ Path: filepath.Join(dir, "css/style.css"), Result: Copied, },
			{ Path: filepath.Join(dir, "myshowcase/hugo-with-nanobox.png"), Result: Overwrited, },
		},
	}

	// Construct LogJSON
	InitLogJSON(cd, sd, dir)
	for _, want := range wantFJ.ContentFiles {
		AddContentFile(want.Path, want.Result)
	}
	for _, want := range wantFJ.NotContentFiles {
		AddNotContentFile(want.Path, want.Result)
	}

	// Save LogJSON
	fn := GetLogFileName()
	outfn := filepath.Join("testdata", fn)
	if err := SaveLogJSON(outfn); err != nil {
		log.Fatalln(err)
	}

	// 
	wantJson := GetLogJSON()

	// Read outfn
	contentBytes, err := os.ReadFile(outfn)
	if err != nil {
		log.Fatalln(err)
	}
	var gotJson LogJSON
	if err := json.Unmarshal(contentBytes, &gotJson); err != nil {
		log.Fatalln(err)
	}

	t.Run("SaveLogJSON", func(t *testing.T) {
		if !reflect.DeepEqual(&gotJson, wantJson) {
			t.Errorf("Compare JSON:\ngot %v\nwant %v", gotJson, wantJson)
		}
	})
}