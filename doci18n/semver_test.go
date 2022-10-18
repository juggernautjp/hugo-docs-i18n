/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Semver function test
*/

package doci18n

import (
	// "os"
	"fmt"
	// "regexp"
	// "io/fs"
	"reflect"
	"testing"
)

// Test for NewSemver()
func TestNewSemver(t *testing.T) {
	t.Run("TestNewSemver()", func(t *testing.T) {
		sv := NewSemver()
		if sv.Major != 0 {
			t.Errorf("Major != 0: %v\n", sv.Major)
		}
		if sv.Minor != 1 {
			t.Errorf("Minor != 1: %v\n", sv.Minor)
		}
		if sv.Patch != 0 {
			t.Errorf("Patch != 0: %v\n", sv.Patch)
		}
		if sv.Prerelease != "" {
			t.Errorf("Patch != \"\": %s\n", sv.Prerelease)
		}
		if sv.Metabuild != "" {
			t.Errorf("Patch != \"\": %s\n", sv.Metabuild)
		}
	})
}

// Test for SetValue()
func TestSetValue(t *testing.T) {
	wantSV := &Semver{
		Major: 3,
		Minor: 4,
		Patch: 5,
		Prerelease: "release",
		Metabuild: "build",
	}
	t.Run("TestSetValue()", func(t *testing.T) {
		gotSV := NewSemver()
		gotSV.SetValue(wantSV.Major, wantSV.Minor, wantSV.Patch, wantSV.Prerelease, wantSV.Metabuild)
		if !reflect.DeepEqual(gotSV, wantSV) {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for IncPatch()
func TestIncPatch(t *testing.T) {
	wantSV := &Semver{
		Major: 0,
		Minor: 1,
		Patch: 1,
		Prerelease: "",
		Metabuild: "",
	}
	t.Run("TestIncPatch()", func(t *testing.T) {
		gotSV := NewSemver()
		gotSV.IncPatch()
		if !reflect.DeepEqual(gotSV, wantSV) {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for IncMinor()
func TestIncMinor(t *testing.T) {
	wantSV := &Semver{
		Major: 0,
		Minor: 2,
		Patch: 0,
		Prerelease: "",
		Metabuild: "",
	}
	t.Run("TestIncMinor()", func(t *testing.T) {
		gotSV := NewSemver()
		gotSV.IncMinor()
		if !reflect.DeepEqual(gotSV, wantSV) {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for IncMajor()
func TestIncMajor(t *testing.T) {
	wantSV := &Semver{
		Major: 1,
		Minor: 0,
		Patch: 0,
		Prerelease: "",
		Metabuild: "",
	}
	t.Run("TestIncMajor()", func(t *testing.T) {
		gotSV := NewSemver()
		gotSV.IncMajor()
		if !reflect.DeepEqual(gotSV, wantSV) {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for SetPrerelease()
func TestSetPrerelease(t *testing.T) {
	wantSV := &Semver{
		Major: 0,
		Minor: 1,
		Patch: 0,
		Prerelease: "release",
		Metabuild: "",
	}
	t.Run("TestSetPrerelease()", func(t *testing.T) {
		gotSV := NewSemver()
		gotSV.SetPrerelease(wantSV.Prerelease)
		if !reflect.DeepEqual(gotSV, wantSV) {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for SetMetabuild()
func TestSetMetabuild(t *testing.T) {
	wantSV := &Semver{
		Major: 0,
		Minor: 1,
		Patch: 0,
		Prerelease: "",
		Metabuild: "metabuild",
	}
	t.Run("TestSetMetabuild()", func(t *testing.T) {
		gotSV := NewSemver()
		gotSV.SetMetabuild(wantSV.Metabuild)
		if !reflect.DeepEqual(gotSV, wantSV) {
			t.Errorf("Compare Semver:\ngot  %v\nwant %v", gotSV, wantSV)
		}
	})
}

// Test for String()
func TestString(t *testing.T) {
	testCases := []struct {
		major, minor, patch int
		prerelease, metabuild, stringify string
	}{
		{ major: 0, minor: 1, patch: 2, prerelease: "", metabuild: "", stringify: "0.1.2"},
		{ major: 0, minor: 1, patch: 2, prerelease: "release", metabuild: "", stringify: "0.1.2-release" },
		{ major: 0, minor: 1, patch: 2, prerelease: "", metabuild: "metabuild", stringify: "0.1.2+metabuild" },
		{ major: 0, minor: 1, patch: 2, prerelease: "release", metabuild: "metabuild", stringify: "0.1.2-release+metabuild" },
	}
	t.Run("TestString()", func(t *testing.T) {
		gotSV := NewSemver()
		for _, want := range testCases {
			gotSV.SetValue(want.major, want.minor, want.patch, want.prerelease, want.metabuild)
			if gotStr := gotSV.String(); gotStr != want.stringify {
				t.Errorf("Compare string:\ngot  %v\nwant %v", gotStr, want.stringify)
			}
		}
	})
}

// Test for ForceValue()
func TestForceValue(t *testing.T) {
	testCases := []struct {
		major, minor, patch int
		prerelease, metabuild, stringify string
		err error
	}{
		{ major: 0, minor: 1, patch: 2, prerelease: "", metabuild: "", stringify: "0.1.2", err: nil},
		{ major: 0, minor: 1, patch: 2, prerelease: "release", metabuild: "", stringify: "0.1.2-release", err: nil },
		{ major: 0, minor: 1, patch: 2, prerelease: "", metabuild: "metabuild", stringify: "0.1.2+metabuild", err: nil },
		{ major: 0, minor: 1, patch: 2, prerelease: "release", metabuild: "metabuild", stringify: "0.1.2-release+metabuild", err: nil },
		{ major: 0, minor: 1, patch: 0, prerelease: "", metabuild: "", stringify: "1.2.3-0123", err: fmt.Errorf("String is not correct SemVer format") },
		{ major: 0, minor: 1, patch: 0, prerelease: "", metabuild: "", stringify: "alpha.beta.1", err: fmt.Errorf("String is not correct SemVer format") },
		{ major: 0, minor: 1, patch: 0, prerelease: "", metabuild: "", stringify: "1.0.0-alpha_beta", err: fmt.Errorf("String is not correct SemVer format") },
		{ major: 0, minor: 1, patch: 0, prerelease: "", metabuild: "", stringify: "9.8.7+meta+meta", err: fmt.Errorf("String is not correct SemVer format") },
	}
	t.Run("TestSetValue()", func(t *testing.T) {
		wantSV := NewSemver() 
		gotSV := NewSemver()
		for _, want := range testCases {
			wantSV.SetValue(want.major, want.minor, want.patch, want.prerelease, want.metabuild)
			err := gotSV.ForceValue(want.stringify)
			if err == nil && err == want.err && !reflect.DeepEqual(gotSV, wantSV) {
				t.Errorf("Compare string: %s\ngot %v\nwant %v", want.stringify, gotSV, wantSV)
			} else if !reflect.DeepEqual(err, want.err) {
				t.Errorf("Not equal error: %s\ngot %s\nwant %s", want.stringify, err, want.err)
			}
		}
	})
}
