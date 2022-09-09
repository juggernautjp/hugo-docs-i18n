/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Copy Content File
*/

package doci18n

import (
	"fmt"
	"os"
	"io/fs"
	// "encoding/json"
)

// Mkdir if the directory dose not exist
func SafeMkdir(path string) error {
	if IsExist(path) {
		return nil
	}
	return os.Mkdir(path, 0755)
}

// Walk directory, while copying the draft file to
func CopyContentDir(srcdir, dstdir string) (PathJSON, error) {
	var pathJSON PathJSON
	var pathPair []PathPair
	SafeMkdir(dstdir)
	// Run WalkDir3()
	err := WalkDir3(srcdir, dstdir, func(srcpath, dstpath string, d fs.DirEntry) error {
		if d.IsDir() {
			if err := SafeMkdir(dstpath); err != nil {
				return err
			}
	 	} else {
			// Copy not content file
			if CopyNotContentFile(srcpath, dstpath) {
				return nil;
			}
			// Copy content file
			b, err := CopyContentFile(srcpath, dstpath)
			if err != nil {
				return err
			}
			// List content file with draft/not-draft flag
			pathPair = append(pathPair, PathPair{Path: srcpath, Draft: b})
		}
		return nil;
	})
 	if err != nil {
		return pathJSON, fmt.Errorf("CopyContentDir: %v", err)
	}
	pathJSON.Files = pathPair
	return pathJSON, nil
}
