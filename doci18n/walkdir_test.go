/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n WalkDir Test
*/

package doci18n

import (
	// "fmt"
	"io/fs"
	"os"
	"log"
	"strings"
	"path/filepath"
	"reflect"
	"testing"
)

// createEmptyFile() create a file contain "content"
/*
func createEmptyFile(file string) {
	if err := os.WriteFile(file, []byte("content"), 0666); err != nil {
		log.Fatal(err)
	}
}
*/

// Test for WalkDir()
func TestWalkDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	os.MkdirAll(filepath.Join(dir, "dir_b"), 0755)
	os.MkdirAll(filepath.Join(dir, "dir_c"), 0755)

	createEmptyFile(filepath.Join(dir, "file_a"))
	createEmptyFile(filepath.Join(dir, "dir_b", "file_b"))
	createEmptyFile(filepath.Join(dir, "dir_c", "file_c"))

	got, err := WalkDir(dir)
	if err != nil {
		t.Errorf("WalkDir: %v", err)
	}
	want := []string{
		filepath.Join(dir, "dir_b", "file_b"),
		filepath.Join(dir, "dir_c", "file_c"),
		filepath.Join(dir, "file_a"),
	}
	t.Run("TestWalkDir()", func(t *testing.T) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("WalkDir():\ngot  %v\nwant %v", got, want)
		}
	})
}

// Test for WalkDir2()
func TestWalkDir2(t *testing.T) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	os.MkdirAll(filepath.Join(dir, "dir_b"), 0755)
	os.MkdirAll(filepath.Join(dir, "dir_c"), 0755)

	createNotEmptyFile(filepath.Join(dir, "file_a"))
	createNotEmptyFile(filepath.Join(dir, "dir_b", "file_b"))
	createNotEmptyFile(filepath.Join(dir, "dir_c", "file_c"))

	wantDir := []string{
		filepath.Join(dir, "dir_b"),
		filepath.Join(dir, "dir_c"),
	}
	wantFile := []string{
		filepath.Join(dir, "dir_b", "file_b"),
		filepath.Join(dir, "dir_c", "file_c"),
		filepath.Join(dir, "file_a"),
	}
	var gotDir, gotFile []string

	t.Run("TestWalkDir2()", func(t *testing.T) {
		_, err := WalkDir2(dir, func(path string, d fs.DirEntry) error {
			if d.IsDir() {
				gotDir = append(gotDir, path)
			} else {
				gotFile = append(gotFile, path)
			}
			return nil;
		})
		if err != nil {
			t.Errorf("TestWalkDir2: %v", err)
		}
		if !reflect.DeepEqual(gotDir, wantDir) {
			t.Errorf("gotDir:\ngot  %v\nwant %v", gotDir, wantDir)
		}
		if !reflect.DeepEqual(gotFile, wantFile) {
			t.Errorf("gotFile:\ngot  %v\nwant %v", gotFile, wantFile)
		}
	})
}

// Test for WalkDir3()
func TestWalkDir3(t *testing.T) {
	const path = "testdata/in"
	const path2 = "testdata/out"
	// Initialize path2 directory
	if IsExist(path2) {
		os.RemoveAll(path2)
		os.Mkdir(path2, 0755)
	}

	t.Run("TestWalkDir3()", func(t *testing.T) {
		// WalkDir3()
		err := WalkDir3(path, path2, func(path, path2 string, d fs.DirEntry) error {
			if d.IsDir() {
				os.Mkdir(path2, 0755)
			} else {
				createEmptyFile(path2)
			}
			return nil;
		})
		if err != nil {
			t.Errorf("TestWalkDir3: %v", err)
		}

		// get the list of directory/file
		want, err := WalkDir(path)
		if err != nil {
			t.Errorf("TestWalkDir3: %v", err)
		}
		got, err := WalkDir(path2)
		if err != nil {
			t.Errorf("TestWalkDir3: %v", err)
		}

		// Test each Field of Struct
		failFunc := func(in, exp, act string) {
			t.Errorf("Path not equal: `%s`\n"+
				"expected: %v\n"+
				"actual  : %v\n", in, exp, act)
		}

		// check if the list is equal
		for k, v := range want {
			srcF := strings.ReplaceAll(v, filepath.FromSlash(path), "")
			dstF := strings.ReplaceAll(got[k], filepath.FromSlash(path2), "")
			if srcF != dstF {
				failFunc(v, srcF, dstF)
			}
		}
		/* can not use reflect.DeepEqual(), 
	  	 because it can not compare "testdata\in\dir_c\file_c" and "testdata\out\dir_c\file_c"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("WalkDir3():\ngot  %v\nwant %v", got, want)
		}
		*/
	})
}
