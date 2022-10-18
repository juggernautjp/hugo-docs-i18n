/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Version function test
*/

package doci18n

import (
	// "os"
	"fmt"
	// "regexp"
	// "io/fs"
	// "reflect"
	"testing"
)

// Test for SaveVersionInfo()
func TestSaveVersionInfo(t *testing.T) {
	t.Run("TestSetValue()", func(t *testing.T) {
		LoadVersionInfo()
		wantSV := GetSemver()
		wantVM := GetVermsg()
		SaveVersionInfo()
		LoadVersionInfo()
		gotSV := GetSemver()
		gotVM := GetVermsg()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
		if gotVM != wantVM {
			t.Errorf("Compare Vermsg:\ngot  %v\nwant %v", gotVM, wantVM)
		}
	})
}

// Test for LoadVersionInfo()
func TestLoadVersionInfo(t *testing.T) {
	t.Run("TestNewSemver()", func(t *testing.T) {
		LoadVersionInfo()
		wantSV := "0.1.0"
		wantVM := DevVersion
		gotSV := GetSemver()
		gotVM := GetVermsg()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
		if gotVM != wantVM {
			t.Errorf("Compare Vermsg:\ngot  %v\nwant %v", gotVM, wantVM)
		}
	})
}

// Test for IncPatch()
func TestIncPatch2(t *testing.T) {
	wantSV := "0.1.1"
	t.Run("TestIncPatch()", func(t *testing.T) {
		LoadVersionInfo()
		IncPatch()
		gotSV := GetSemver()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for IncMinor()
func TestIncMinor2(t *testing.T) {
	wantSV := "0.2.0"
	t.Run("TestIncMinor()", func(t *testing.T) {
		LoadVersionInfo()
		IncMinor()
		gotSV := GetSemver()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for IncMajor()
func TestIncMajor2(t *testing.T) {
	wantSV := "1.0.0"
	wantVM := ReleaseVersion
	t.Run("TestIncMajor()", func(t *testing.T) {
		LoadVersionInfo()
		IncMajor()
		gotSV := GetSemver()
		gotVM := GetVermsg()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
		if gotVM != wantVM {
			t.Errorf("Compare Vermsg:\ngot  %v\nwant %v", gotVM, wantVM)
		}
	})
}

// Test for SetPrerelease()
func TestSetPrerelease2(t *testing.T) {
	prerelease := "release"
	wantSV := fmt.Sprintf("0.1.0-%s", prerelease)
	t.Run("TestSetPrerelease()", func(t *testing.T) {
		LoadVersionInfo()
		SetPrerelease(prerelease)
		gotSV := GetSemver()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for SetMetabuild()
func TestSetMetabuild2(t *testing.T) {
	metabuild := "metabuild"
	wantSV := fmt.Sprintf("0.1.0+%s", metabuild)
	t.Run("TestSetMetabuild()", func(t *testing.T) {
		LoadVersionInfo()
		SetMetabuild(metabuild)
		gotSV := GetSemver()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for SetSemver()
func TestSetSemver(t *testing.T) {
	wantSV := "3.4.5-release+build"
	t.Run("TestSetValue()", func(t *testing.T) {
		LoadVersionInfo()
		if err := SetSemver(wantSV); err != nil {
			t.Errorf("Error SetSemver: %s\n", wantSV)
		}
		gotSV := GetSemver()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for SetVermsg()
func TestSetVermsg(t *testing.T) {
	wantVM := "Beta version"
	t.Run("TestSetValue()", func(t *testing.T) {
		LoadVersionInfo()
		SetVermsg(wantVM)
		gotVM := GetVermsg()
		if gotVM != wantVM {
			t.Errorf("Compare Vermsg:\ngot  %v\nwant %v", gotVM, wantVM)
		}
	})
}

