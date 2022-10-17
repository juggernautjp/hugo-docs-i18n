/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Copy Content File
*/

package doci18n

import (
	// "fmt"
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
func CopyContentDir(srcdir, dstdir, outfn string) (PathJSON, error) {
	var pathJSON PathJSON
	var pathPair []PathPair
	InitLogJSON(GetCurrentString(), srcdir, dstdir)
	SafeMkdir(dstdir)
	// Run WalkDir3()
	err := WalkDir3(srcdir, dstdir, func(srcpath, dstpath string, d fs.DirEntry) error {
		if d.IsDir() {
			if err := SafeMkdir(dstpath); err != nil {
				return err
			}
	 	} else {
			isExist := false
			if IsExist(dstpath) {
				isExist = true
			}
			if ret, err := CopyNotContentFile(srcpath, dstpath); ret {
				// Copy not-content file
				if err != nil {
					AddNotContentFile(dstpath, Failed)
					return err
				} else if isExist {
					// Overwrite
					AddNotContentFile(dstpath, Overwrited)
				} else {
					AddNotContentFile(dstpath, Copied)
				}
				return nil;
			}
			// Copy content file
			b, err := CopyContentFile(srcpath, dstpath)
			if err != nil {
				AddContentFile(dstpath, Failed)
				return err
			} else if isExist {
				// Do not overwrite
			} else {
				// Copy
				AddContentFile(dstpath, Copied)
			}
			// List content file with draft/not-draft flag
			pathPair = append(pathPair, PathPair{Path: srcpath, Draft: b})
		}
		return nil;
	})
 	if err != nil {
		// return pathJSON, fmt.Errorf("CopyContentDir: %v", err)
		return pathJSON, err
	}
	pathJSON.Files = pathPair
	// Save JSON
	if err := SaveLogJSON(outfn); err != nil {
		return pathJSON, err
	}
	return pathJSON, nil
}
