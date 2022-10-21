/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Version JSON data
*/

package doci18n

import (
	"os"
	"fmt"
	// "regexp"
	// "strconv"
	// "io/fs"
	"encoding/json"
	// "log"
	_ "embed"
)

//go:embed version.json
var contentBytes []byte

// Global variables
var GlobalVersionInfo VersionInfo
const DevVersion = "Dev version"
const ReleaseVersion = "Release version"
const VersionFile = "version.json"

// VersionInfo data
type VersionInfo struct {
	 SemVer string `json:"semver"`
	 VerMsg string `json:"vermsg"`
}

// Can Set Version Info
func CanSetVersionInfo() bool {
	return IsExist(VersionFile)
}

// Initialize Version Info
func InitVersionInfo() {
	sv := NewSemver()
	GlobalVersionInfo.SemVer = sv.String()
	GlobalVersionInfo.VerMsg = DevVersion
}

// Load VersionInfo data from JSON
func LoadVersionInfo() error {
	// Emberd need "version.json" file
	if !IsExist(VersionFile) {
		sv := NewSemver()
		GlobalVersionInfo.SemVer = sv.String()
		GlobalVersionInfo.VerMsg = DevVersion
		return nil
	}
	// load JSON
	contentBytes, err := os.ReadFile(VersionFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(contentBytes, &GlobalVersionInfo); err != nil {
		return err
	}
	return nil
}

// Load VersionInfo data from JSON
func GetVersionInfo() error {
	if !IsExist(VersionFile) {
		return fmt.Errorf("not found VersionFile")
	}
	if err := json.Unmarshal(contentBytes, &GlobalVersionInfo); err != nil {
		return err
	}
	return nil
}

// Load VersionInfo data from JSON
func GetSemver() string {
	return GlobalVersionInfo.SemVer
}

func GetVermsg() string {
	return GlobalVersionInfo.VerMsg
}

// Inc SemVer Info
func IncMajor() error {
	svStr := GlobalVersionInfo.SemVer
	sv := NewSemver()
	if err := sv.ForceValue(svStr); err != nil {
		return err
	}
	if sv.Major == 0 {
		GlobalVersionInfo.VerMsg = ReleaseVersion
	}
	sv.IncMajor()
	GlobalVersionInfo.SemVer = sv.String()
	return nil
}

func IncMinor() error {
	svStr := GlobalVersionInfo.SemVer
	sv := NewSemver()
	if err := sv.ForceValue(svStr); err != nil {
		return err
	}
	sv.IncMinor()
	GlobalVersionInfo.SemVer = sv.String()
	return nil
}

func IncPatch() error {
	svStr := GlobalVersionInfo.SemVer
	sv := NewSemver()
	if err := sv.ForceValue(svStr); err != nil {
		return err
	}
	sv.IncPatch()
	GlobalVersionInfo.SemVer = sv.String()
	return nil
}

func SetPrerelease(m string) error {
	svStr := GlobalVersionInfo.SemVer
	sv := NewSemver()
	if err := sv.ForceValue(svStr); err != nil {
		return err
	}
	sv.SetPrerelease(m)
	GlobalVersionInfo.SemVer = sv.String()
	return nil
}

func SetMetabuild(m string) error {
	svStr := GlobalVersionInfo.SemVer
	sv := NewSemver()
	if err := sv.ForceValue(svStr); err != nil {
		return err
	}
	sv.SetMetabuild(m)
	GlobalVersionInfo.SemVer = sv.String()
	return nil
}

func SetSemver(str string) error {
	// svStr := GlobalVersionInfo.SemVer
	sv := NewSemver()
	if err := sv.ForceValue(str); err != nil {
		return err
	}
	GlobalVersionInfo.SemVer = sv.String()
	return nil
}

func SetVermsg(str string) {
	GlobalVersionInfo.VerMsg = str
}

// Save JSON data to file
func SaveVersionInfo() error {
	// open the file
	outfile, err := os.Create(VersionFile)
	if err != nil {
		return err
	}
	// close the file
	defer outfile.Close()

	// write DraftJSON
	sd, err := json.MarshalIndent(&GlobalVersionInfo, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(outfile, string(sd))
	return nil
}
