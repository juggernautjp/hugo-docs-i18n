/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n CountDraftFile Test
*/

package doci18n

import (
	// "fmt"
	// "io/fs"
	"os"
	// "log"
	// "strings"
	"path/filepath"
	"encoding/json"
	"reflect"
	"testing"
)


// Test for CountDraftFile()
func TestCountDraftFile(t *testing.T) {
	dir := "testdata/zn"
	wantPJ := PathJSON{
		Files: []PathPair{
			{ Path: filepath.Join(dir, "dir_a/draft_a.md"), Title: "title draft_a", Draft: true, },
			{ Path: filepath.Join(dir, "dir_a/notdraft_a.md"), Title: "title notdraft_a", Draft: false, },
			{ Path: filepath.Join(dir, "dir_b/draft_b.md"), Title: "title draft_b", Draft: true, },
			{ Path: filepath.Join(dir, "dir_b/notdraft_b.md"), Title: "title notdraft_b", Draft: false, },
			{ Path: filepath.Join(dir, "draft_c.md"), Title: "title draft_c", Draft: true, },
			{ Path: filepath.Join(dir, "myshowcase/no_draft.md"), Title: "YAML NoDraft Markdown", Draft: false, },
			{ Path: filepath.Join(dir, "myshowcase/no_frontmatter.md"), Title: "", Draft: false, },
			{ Path: filepath.Join(dir, "notdraft_c.md"), Title: "title notdraft_c", Draft: false, },
		},
	}
	wantCD := CountDraft{
		Draft: 3,
		NotDraft: 5,
		Total: 5+3,
		Ratio: float32(5)/float32(5+3),
	}

	// Run CountDraftFile()
	gotPJ, gotCD, err := CountDraftFile(dir)
	if err != nil {
		t.Errorf("TestCountDraftFile: %v", err)
	}

	// Case 1: compare JSON
	t.Run("Case 1: Compare JSON", func(t *testing.T) {
		if !reflect.DeepEqual(gotPJ, wantPJ) {
			t.Errorf("Compare JSON:\ngot  %v\nwant %v", gotPJ, wantPJ)
		}
	})

	// Case 2: compare CountDraft
	t.Run("Case 2: Compare CountDraft", func(t *testing.T) {
		if !reflect.DeepEqual(gotCD, wantCD) {
			t.Errorf("Compare JSON:\ngot  %v\nwant %v", gotCD, wantCD)
		}
	})
}

// Test for CountDraftFile()
func TestSaveCountJSONFile(t *testing.T) {
	dir := "testdata/zn"
	fn := filepath.Join("testdata", "zn.json")

	var wantJson DraftJSON = DraftJSON{
		Count: CountDraft{
			Draft: 3,
			NotDraft: 5,
			Total: 5+3,
			Ratio: float32(5)/float32(5+3),
			},
		Paths: PathJSON{
			Files: []PathPair{
				{ Path: filepath.Join(dir, "dir_a/draft_a.md"), Title: "title draft_a", Draft: true, },
				{ Path: filepath.Join(dir, "dir_a/notdraft_a.md"), Title: "title notdraft_a", Draft: false, },
				{ Path: filepath.Join(dir, "dir_b/draft_b.md"), Title: "title draft_b", Draft: true, },
				{ Path: filepath.Join(dir, "dir_b/notdraft_b.md"), Title: "title notdraft_b", Draft: false, },
				{ Path: filepath.Join(dir, "draft_c.md"), Title: "title draft_c", Draft: true, },
				{ Path: filepath.Join(dir, "myshowcase/no_draft.md"), Title: "YAML NoDraft Markdown", Draft: false, },
				{ Path: filepath.Join(dir, "myshowcase/no_frontmatter.md"), Title: "", Draft: false, },
				{ Path: filepath.Join(dir, "notdraft_c.md"), Title: "title notdraft_c", Draft: false, },
			},
		},
	}

	// Run CountDraftFile()
	gotPJ, gotCD, err := CountDraftFile(dir)
	if err != nil {
		t.Errorf("TestSaveCountJSONFile: %v", err)
	}

	// Run SaveCountJSONFile()
	if err := SaveCountJSONFile(fn, gotPJ, gotCD); err != nil {
		t.Errorf("TestSaveCountJSONFile: %v", err)
	}

	// Read from saved JSON file
	var gotJson DraftJSON
	cb, err := os.ReadFile(fn)
	if err != nil {
		t.Errorf("TestSaveCountJSONFile: %v", err)
	}
	if err := json.Unmarshal(cb, &gotJson); err != nil {
		t.Errorf("TestSaveCountJSONFile: %v", err)
	}

	// Case 3: Test SaveCountJSONFile()
	t.Run("Case 3: SaveCountJSONFile", func(t *testing.T) {
		if !reflect.DeepEqual(gotJson, wantJson) {
			t.Errorf("Compare JSON:\ngot  %v\nwant %v", gotJson, wantJson)
		}
	})
}