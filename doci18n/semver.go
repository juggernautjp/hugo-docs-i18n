/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Semver function set
regstr value is cite from: https://regex101.com/r/vkijKf/1/
*/

package doci18n

import (
	// "os"
	"fmt"
	"regexp"
	"strconv"
	// "io/fs"
)

// semver regexp
const regstr = `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

// Data Semser
type Semver struct {
	Major int
	Minor int
	Patch int
	Prerelease string
	Metabuild string
}

// Increment Semver.Major + 1
func (sv *Semver) IncMajor() {
	sv.Major++
	sv.Minor = 0
	sv.Patch = 0
}

// Increment Semver.Minor + 1
func (sv *Semver) IncMinor() {
	sv.Minor++
	sv.Patch = 0
}

// Increment Semver.Patch + 1
func (sv *Semver) IncPatch() {
	sv.Patch++
}

// Set string to Semver.Prerelease
func (sv *Semver) SetPrerelease(m string) {
	sv.Prerelease = m
}

// Set string to Semver.Metabuild
func (sv *Semver) SetMetabuild(m string) {
	sv.Metabuild = m
}

// Force to set the value of Semver
func (sv *Semver) SetValue(ma, mi, pa int, pr, mb string) {
	sv.Major = ma
	sv.Minor = mi
	sv.Patch = pa
	sv.Prerelease = pr
	sv.Metabuild = mb
}

// Transform string value to struct Semver
func (sv *Semver) String() string {
	if sv.Prerelease != "" && sv.Metabuild != "" {
		return fmt.Sprintf("%d.%d.%d-%s+%s", sv.Major, sv.Minor, sv.Patch, sv.Prerelease, sv.Metabuild)
	} else if sv.Prerelease != "" {
		return fmt.Sprintf("%d.%d.%d-%s", sv.Major, sv.Minor, sv.Patch, sv.Prerelease)
	} else if sv.Metabuild != "" {
		return fmt.Sprintf("%d.%d.%d+%s", sv.Major, sv.Minor, sv.Patch, sv.Metabuild)
	} else {
		return fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
	}
}

// Constructor of Semver
func NewSemver() *Semver {
	return &Semver{Major: 0, Minor: 1, Patch: 0, Prerelease: "", Metabuild: ""}
}

// Force to set the value to struct Semver with SemVer format 
func (sv *Semver) ForceValue(str string) error {
	regex := regexp.MustCompile(regstr)
	sub := regex.FindStringSubmatch(str)
	if sub == nil {
		return fmt.Errorf("String is not correct SemVer format")
	}
	// regex parse error
	if sub[1] == "" {
		return fmt.Errorf("Major is not digit")
	}
	if sub[2] == "" {
		return fmt.Errorf("Minor is not digit")
	}
	if sub[3] == "" {
		return fmt.Errorf("Patch is not digit")
	}
	// Parse Semver
	var major, minor, patch int64
	var err error
	major, err = strconv.ParseInt(sub[1], 10, 0)
	if err != nil {
		return fmt.Errorf("ParseInt sub[1] error: %s", sub[1])
	}
	minor, err = strconv.ParseInt(sub[2], 10, 0)
	if err != nil {
		return fmt.Errorf("ParseInt sub[2] error: %s", sub[2])
	}
	patch, err = strconv.ParseInt(sub[3], 10, 0)
	if err != nil {
		return fmt.Errorf("ParseInt sub[3] error: %s", sub[3])
	}
	// Set value
	sv.Major = int(major)
	sv.Minor = int(minor)
	sv.Patch = int(patch)
	sv.Prerelease = sub[4]
	sv.Metabuild = sub[5]
	return nil
}
