/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n WalkDir
*/

package doci18n

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func WalkDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %w", err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			// Recursively calls WalkDir in the case of a directory
			p, err := WalkDir(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("WalkDir %s: %w", filepath.Join(dir, file.Name()), err)
			}
			// Merge into the caller's "paths" variable.
			paths = append(paths, p...)
		} else {
			// Now that we've reached a leaf (file) in the directory tree, we'll add it to "paths" variable.
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}
	fmt.Println(paths)
	return paths, nil
}

func WalkDir2(dir string, walkDirFunc func(path string, d fs.DirEntry) error) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %w", err)
	}

	var paths []string
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			// call walkDirFunc() for directory `path`
			err1 := walkDirFunc(path, file)
			if err1 != nil {
				return nil, fmt.Errorf("walkDirFunc %s: %w", path, err1)
			}
			// Recursively calls WalkDir in the case of a directory
			p, err := WalkDir2(path, walkDirFunc)
			if err != nil {
				return nil, fmt.Errorf("WalkDir2 %s: %w", path, err)
			}
			// Merge into the caller's "paths" variable.
			paths = append(paths, p...)
		} else {
			// Now that we've reached a leaf (file) in the directory tree, we'll add it to "paths" variable.
			paths = append(paths, path)
		}
	}
	// fmt.Println(paths)
	return paths, nil
}
