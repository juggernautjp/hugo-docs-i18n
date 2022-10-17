/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Count Draft File
*/

package doci18n

import (
	"os"
	"fmt"
	"io/fs"
	"encoding/json"

	"github.com/gohugoio/hugo/hugofs/files"
)

// DraftJSON
type DraftJSON struct {
	Count CountDraft `json:"count"`
	Paths PathJSON `json:"paths"`
}

// Hugo file (path/value pair) data
type PathPair struct {
	Path string `json:"path"`
	Title string `json:"title"`
	Draft bool `json:"draft"`
}

type PathJSON struct {
	Files []PathPair `json:"files"`
}

// Count Not-Draft(translated)/Draft(not-translated) file
type CountDraft struct {
	Draft int `json:"draft"`
	NotDraft int `json:"notdraft"`
	Total int `json:"total"`
	Ratio float32 `json:"ratio"`
}

func (cd *CountDraft) Sum() int {
	cd.Total = cd.Draft + cd.NotDraft
	if cd.Total == 0 {
		cd.Ratio = 0
	} else {
		cd.Ratio = float32(cd.NotDraft) / float32(cd.Total)
	}
	return cd.Total
}

func (cd *CountDraft) Add(b bool) {
	if b {
		cd.Draft++
	} else {
		cd.NotDraft++
	}
}

// Walk directory, while checking if the file is draft or not?
func CountDraftFile(dir string) (PathJSON, CountDraft, error) {
	var pathJSON PathJSON
	var pathPair []PathPair
	var cd CountDraft
	_, err := WalkDir2(dir, func(path string, d fs.DirEntry) error {
		if d.IsDir() {
			return nil;
	 	} else {
			// when file is not content file,
			if !IsRegularFile(path) || !files.IsContentFile(path) {
				return nil
			}
			// when file is content file, check if it is draft or not?
			b, s, err := IsDraftFile(path)
			if err != nil {
				return err
			}
			pathPair = append(pathPair, PathPair{Path: path, Title: s, Draft: b})
			cd.Add(b)
		}
		return nil;
	})
 	if err != nil {
		return pathJSON, cd, fmt.Errorf("CountDraftFile: %v", err)
	}
	pathJSON.Files = pathPair
	cd.Sum()
	return pathJSON, cd, nil
}

// Save JSON data to file
func SaveCountJSONFile(outfn string, pj PathJSON, cd CountDraft) error {
	var dj DraftJSON
	dj.Count = cd
	dj.Paths = pj

	// open the file
	outfile, err := os.Create(outfn)
	if err != nil {
		// return fmt.Errorf("Error when opening file: %s", err)
		return err
	}
	// close the file
	defer outfile.Close()

	// write DraftJSON
	sd, err := json.MarshalIndent(dj, "", "  ")
	if err != nil {
		// return fmt.Errorf("Error when marshaling DraftJSON: %s", err)
		return err
	}
	fmt.Fprintln(outfile, string(sd))
	return nil
}