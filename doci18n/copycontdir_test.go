/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n CopyContentDir Test

TODO: 最初のコピー (コピー先が空の場合) の後、元 (英語) のコンテンツが変更された場合の反映
*/

package doci18n

import (
	// "fmt"
	// "io/fs"
	"os"
	// "log"
	// "strings"
	"path/filepath"
	"reflect"
	"testing"
)

// Test for CopyContentDir()
func TestCopyContentDir(t *testing.T) {
	// Prepare src directory
	srcdir := "testdata/zn"
	wantPJ := PathJSON{
		Files: []PathPair{
			{ Path: filepath.Join(srcdir, "dir_a/draft_a.md"), Draft: false, },
			{ Path: filepath.Join(srcdir, "dir_a/notdraft_a.md"), Draft: true, },
			{ Path: filepath.Join(srcdir, "dir_b/draft_b.md"), Draft: false, },
			{ Path: filepath.Join(srcdir, "dir_b/notdraft_b.md"), Draft: true, },
			{ Path: filepath.Join(srcdir, "draft_c.md"), Draft: false, },
			{ Path: filepath.Join(srcdir, "notdraft_c.md"), Draft: true, },
		},
	}

	// Prepare dst directory
	dstdir := "testdata/ko"
	if IsExist(dstdir) {
		os.RemoveAll(dstdir)
		os.Mkdir(dstdir, 0755)
	}
	wantDT := []string {
		filepath.Join(dstdir, "css", "style.css"),
		filepath.Join(dstdir, "dir_a", "notdraft_a.md"),
		filepath.Join(dstdir, "dir_b", "notdraft_b.md"),
		filepath.Join(dstdir, "images", "hugo-with-nanobox.png"),
		filepath.Join(dstdir, "notdraft_c.md"),
	}
	wantDraft := []string {
		"dir_a/notdraft_a.md",
		"dir_b/notdraft_b.md",
		"notdraft_c.md",
	}

	// Run CopyContentDir()
	gotPJ, err := CopyContentDir(srcdir, dstdir)
	if err != nil {
		t.Errorf("TestCopyContentDir: %v", err)
	}

	// Case 1: compare JSON
	t.Run("Case 1: Compare JSON", func(t *testing.T) {
		if !reflect.DeepEqual(gotPJ, wantPJ) {
			t.Errorf("Compare JSON:\ngot  %v\nwant %v", gotPJ, wantPJ)
		}
	})

	// Case 2: compare file/directory tree
	t.Run("Case 2: Compare Tree", func(t *testing.T) {
		gotDT, err := WalkDir(dstdir)
		if err != nil {
			t.Errorf("TestCopyContentDir: %v", err)
		}
		if !reflect.DeepEqual(gotDT, wantDT) {
			t.Errorf("Compare JSON:\ngot  %v\nwant %v", gotDT, wantDT)
		}
	})

	// Case 3: compare CountDraft
	t.Run("Case 3: Check Draft", func(t *testing.T) {
		for _, fn := range wantDraft {
			infname := filepath.Join(dstdir, fn)
			b, _, err := IsDraftFile(infname)
			if err != nil {
				t.Errorf("Error: IsDraftFile: %v", err)
			}
			if !b {
				t.Errorf("Copied file is Not draft: %v", infname)
			}
		}
	})
}
