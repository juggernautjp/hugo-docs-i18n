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

// WalkDir() traverse directory tree
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
	// fmt.Println(paths)
	return paths, nil
}

// WalkDir2() traverse directory tree with walkDirFunc()
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
			if err := walkDirFunc(path, file); err != nil {
				return nil, fmt.Errorf("walkDirFunc %s: %w", path, err)
			}
			// Recursively calls WalkDir in the case of a directory
			p, err := WalkDir2(path, walkDirFunc)
			if err != nil {
				return nil, fmt.Errorf("WalkDir2 %s: %w", path, err)
			}
			// Merge into the caller's "paths" variable.
			paths = append(paths, p...)
		} else {
			// Now that we've reached a leaf (file) in the directory tree

			// call walkDirFunc() for directory `path`
			if err := walkDirFunc(dir, file); err != nil {
				return nil, fmt.Errorf("walkDirFunc %s: %w", path, err)
			}
			// we'll add it to "paths" variable.
			paths = append(paths, path)
		}
	}
	// fmt.Println(paths)
	return paths, nil
}

// WalkDir3() traverse src directory tree with walkDirFunc()
func WalkDir3(srcdir, dstdir string, walkDirFunc func(srcpath, dstpath string, d fs.DirEntry) error) error {
	files, err := os.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("ReadDir: %w", err)
	}

	// var paths []string
	for _, file := range files {
		srcpath := filepath.Join(srcdir, file.Name())
		dstpath := filepath.Join(dstdir, file.Name())
		if file.IsDir() {
			// call walkDirFunc() for directory `path`
			if err := walkDirFunc(srcpath, dstpath, file); err != nil {
				return fmt.Errorf("walkDirFunc %s: %w", srcpath, err)
			}
			// Recursively calls WalkDir in the case of a directory
			if err := WalkDir3(srcpath, dstpath, walkDirFunc); err != nil {
				return fmt.Errorf("WalkDir3 %s: %w", srcpath, err)
			}
			// Merge into the caller's "paths" variable.
			// paths = append(paths, p...)
		} else {
			// Now that we've reached a leaf (file) in the directory tree

			// call walkDirFunc() for directory `path`
			if err := walkDirFunc(srcpath, dstpath, file); err != nil {
				return fmt.Errorf("walkDirFunc %s: %w", srcpath, err)
			}
			// we'll add it to "paths" variable.
			// paths = append(paths, srcpath)
		}
	}
	// return paths, nil
	return nil
}

