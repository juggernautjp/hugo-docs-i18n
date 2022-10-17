/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo log while copy content directory
*/

package doci18n

import (
	"os"
	"fmt"
	// "io/fs"
	"time"
	"encoding/json"
)

// Result code of copy content
const (
	NotExist = iota
	Copied
	Overwrited
	NotUpdated
	Failed
)

// Date/Time Format
// const longForm = time.RFC1123Z
const longForm = "2006-01-02 15:04:05 -0700 MST"
const shortForm = "2006-01-02_15-04-05"

// Global variable
var globalLogJSON *LogJSON

// LogJSON
type LogJSON struct {
	Header SummaryJSON `json:"header"`
	Paths FileJSON `json:"paths"`
}

// Hugo file (path/value pair) data
type CopyFile struct {
	Path string `json:"path"`
	Result int `json:"result"`
}

type FileJSON struct {
	ContentFiles []CopyFile `json:"contentfiles"`
	NotContentFiles []CopyFile `json:"notcontentfiles"`
}

// Summary
type SummaryJSON struct {
	Date string `json:"date"`
	SrcDir string `json:"srcdir"`
	DstDir string `json:"dstdir"`
}

// Get Log file name
func GetLogFileName() string {
	t := time.Now()
	tz, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s.log", t.In(tz).Format(shortForm))
}

// Get Now() string
func GetCurrentString() string {
	t := time.Now()
	tz, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	return t.In(tz).Format(longForm)
}

// Initialize LogJSON
func InitLogJSON(dt, sd, dd string) {
	globalLogJSON = NewLogJSON(dt, sd, dd)
}

func GetLogJSON() *LogJSON {
	return globalLogJSON
}

// Add CopyFile to LogJSON
func AddContentFile(p string, r int) error {
	if globalLogJSON == nil {
		return fmt.Errorf("Error: LogJSON is not initialized.")
	}
	globalLogJSON.AddContentFile(p, r)
	return nil
}

func AddNotContentFile(p string, r int) error {
	if globalLogJSON == nil {
		return fmt.Errorf("Error: LogJSON is not initialized.")
	}
	globalLogJSON.AddNotContentFile(p, r)
	return nil
}

// Save JSON data to file
func SaveLogJSON(outfn string) error {
	if globalLogJSON == nil {
		return fmt.Errorf("Error: LogJSON is not initialized.")
	}
	return globalLogJSON.SaveLogJSON(outfn)
}



// Constructor of LogJSON
func NewLogJSON(dt, sd, dd string) *LogJSON {
	return &LogJSON{ 
		Header: SummaryJSON{ Date: dt, SrcDir: sd, DstDir: dd, },
		Paths: FileJSON{ ContentFiles: []CopyFile{}, NotContentFiles: []CopyFile{}, },
	}
}

// Add CopyFile to LogJSON
func (lg *LogJSON) AddContentFile(p string, r int) {
	lg.Paths.ContentFiles = append(lg.Paths.ContentFiles, CopyFile{ Path: p, Result: r, })
}

func (lg *LogJSON) AddNotContentFile(p string, r int) {
	lg.Paths.NotContentFiles = append(lg.Paths.NotContentFiles, CopyFile{ Path: p, Result: r, })
}

// Save JSON data to file
func (lg *LogJSON) SaveLogJSON(outfn string) error {
	if outfn == "" {
		return fmt.Errorf("output file not specified")
	}

	// open the file
	outfile, err := os.Create(outfn)
	if err != nil {
		// return fmt.Errorf("Error when opening file: %s", err)
		return err
	}
	// close the file
	defer outfile.Close()

	// write DraftJSON
	sd, err := json.MarshalIndent(lg, "", "  ")
	if err != nil {
		// return fmt.Errorf("Error when marshaling DraftJSON: %s", err)
		return err
	}
	fmt.Fprintln(outfile, string(sd))
	return nil
}


