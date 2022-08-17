/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n WalkDir Test
*/

package doci18n

import (
	"fmt"
	"io/fs"
	"os"
	"log"
	"path/filepath"
	"reflect"
	"testing"
)

func createEmptyFile(file string) {
	if err := os.WriteFile(file, []byte("content"), 0666); err != nil {
		log.Fatal(err)
	}
}

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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WalkDir():\ngot  %v\nwant %v", got, want)
	}
}

func TestWalkDir2(t *testing.T) {
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

	got, err := WalkDir2(dir, func(path string, d fs.DirEntry) error {
		if d.IsDir() {
			fmt.Println("walkDirFunc dir: ", path)
		} else {
			fmt.Println("walkDirFunc file: ", path)
		}
		return nil;
	})
	if err != nil {
		t.Errorf("WalkDir2: %v", err)
	}
	want := []string{
		filepath.Join(dir, "dir_b", "file_b"),
		filepath.Join(dir, "dir_c", "file_c"),
		filepath.Join(dir, "file_a"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WalkDir2():\ngot  %v\nwant %v", got, want)
	}
}
