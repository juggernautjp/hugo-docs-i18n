/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Version function test
*/

package doci18n

import (
	"os"
	"fmt"
	"log"
	// "regexp"
	// "io/fs"
	// "reflect"
	"testing"
)

// Test for LoadVersionInfo()
// Test for GetVersionInfo()
func TestGetVersionInfo(t *testing.T) {
	t.Run("TestGetVersionInfo()", func(t *testing.T) {
		GetVersionInfo()
		wantSV := "0.1.0"
		wantVM := "Dev version"
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
	t.Run("TestIncPatch2()", func(t *testing.T) {
		GetVersionInfo()
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
	t.Run("TestIncMinor2()", func(t *testing.T) {
		GetVersionInfo()
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
	t.Run("TestIncMajor2()", func(t *testing.T) {
		// GetVersionInfo()
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
	sv := "1.0.0"
	prerelease := "release"
	wantSV := fmt.Sprintf("%s-%s", sv, prerelease)
	t.Run("TestSetPrerelease2()", func(t *testing.T) {
		// GetVersionInfo()
		SetPrerelease(prerelease)
		gotSV := GetSemver()
		if gotSV != wantSV {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for SetMetabuild()
func TestSetMetabuild2(t *testing.T) {
	sv := "1.0.0-release"
	metabuild := "metabuild"
	wantSV := fmt.Sprintf("%s+%s",sv,  metabuild)
	t.Run("TestSetMetabuild2()", func(t *testing.T) {
		// GetVersionInfo()
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
	t.Run("TestSetSemver()", func(t *testing.T) {
		// GetVersionInfo()
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
	wantVM := DevVersion
	t.Run("TestSetVermsg()", func(t *testing.T) {
		// GetVersionInfo()
		SetVermsg(wantVM)
		gotVM := GetVermsg()
		if gotVM != wantVM {
			t.Errorf("Compare Vermsg:\ngot  %v\nwant %v", gotVM, wantVM)
		}
	})
}

// Test for SaveVersionInfo()
func TestSaveVersionInfo(t *testing.T) {
	t.Run("TestSaveVersionInfo()", func(t *testing.T) {
		// LoadVersionInfo()
		wantSV := "3.4.5-release+build"
		wantVM := DevVersion
		SaveVersionInfo()
		if err := LoadVersionInfo(); err != nil {
			log.Fatalln(err)
		}
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

func TestLoadVersionInfo(t *testing.T) {
	t.Run("TestLoadVersionInfo()", func(t *testing.T) {
		if err := LoadVersionInfo(); err != nil {
			t.Errorf("Error LoadVersionInfo: %s\n", err)
		}
		wantSV := "3.4.5-release+build"
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

// Test for CanSetVersionInfo()
func TestCanSetVersionInfo(t *testing.T) {
	t.Run("TestCanSetVersionInfo(true)", func(t *testing.T) {
		wantVI := true
		gotVI := CanSetVersionInfo()
		if gotVI != gotVI {
			t.Errorf("VersionInfo: \ngot  %v\nwant %v", gotVI, wantVI)
		}
	})
	os.Remove(VersionFile)
	t.Run("TestCanSetVersionInfo(false)", func(t *testing.T) {
		wantVI := false
		gotVI := CanSetVersionInfo()
		if gotVI != gotVI {
			t.Errorf("VersionInfo: \ngot  %v\nwant %v", gotVI, wantVI)
		}
	})
}

// Before/After Test Hook function
func BeforeTestHook() {
	if !CanSetVersionInfo() {
		log.Fatalln("Before version_test.go, you should initialize doci18n/version.json")
	}
}

func AfterTestHook() {
	// Initialize VersionInfor
	InitVersionInfo() 
	// Save VersionInfo JSON file
	SaveVersionInfo()
}

// TestMain
func TestMain(m *testing.M) {
	fmt.Println("前処理")
	BeforeTestHook()
	status := m.Run()
	fmt.Println("後処理")
	AfterTestHook()
	os.Exit(status)
}