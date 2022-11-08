/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n l10n Test
*/

package locale

import (
	// "os"
	// "log"
	// "encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)


// Test for InitTranslation()
func TestInitTranslation(t *testing.T) {
	var wantBundle Bundle = Bundle{
		Pairs: []KeyPair{
			{ Key: "rootCmdShort", Message: "rootCmdShort EnMessage" },
			{ Key: "rootCmdLong", Message: "rootCmdLong EnMessage" },
		},
	}

	t.Run("Before InitTranslation", func(t *testing.T) {
		if !EnMessages.isEmpty() {
			t.Errorf("EnMessage is not empty.\n")
		}
	})

	t.Run("InitTranslation", func(t *testing.T) {
		// Test Run
		err := InitTranslation()
		if err != nil {
			t.Fatalf("Failed to initialize: %s\n", err)
		}
		if !reflect.DeepEqual(EnMessages, wantBundle) {
			t.Errorf("Compare Bundle with JSON:\ngot  %v\nwant %v", EnMessages, wantBundle)
		}
	})

	t.Run("After InitTranslation", func(t *testing.T) {
		if EnMessages.isEmpty() {
			t.Errorf("EnMessage is empty.\n")
		}
	})
}

// Test for func LoadTranslation()
func TestLoadTranslation(t *testing.T) {
	fn := filepath.Join("testdata", "l10n.ja.json")

	var wantBundle Bundle = Bundle{
		Pairs: []KeyPair{
			{ Key: "rootCmdShort", Message: "rootCmdShort JaMessage" },
			{ Key: "rootCmdLong", Message: "rootCmdLong JaMessage" },
		},
	}

	t.Run("LoadTranslation", func(t *testing.T) {
		// Test Run
		err := LoadTranslation(fn)
		if err != nil {
			t.Fatalf("Failed to initialize: %s\n", err)
		}
		if !reflect.DeepEqual(TargetMessages, wantBundle) {
			t.Errorf("Compare Bundle with JSON:\ngot  %v\nwant %v", TargetMessages, wantBundle)
		}
	})
}

// Test for func T()
func TestT(t *testing.T) {
	// Test with en.json
	TargetMessages.Pairs = []KeyPair{}
	var testCases = map[string]string {
		"rootCmdShort": "rootCmdShort EnMessage",
		"rootCmdLong": "rootCmdLong EnMessage",
	}

	t.Run("T with en.json", func(t *testing.T) {
		// Test Run
		err := InitTranslation()
		if err != nil {
			t.Fatalf("Failed to initialize: %s\n", err)
		}
		for k, want := range testCases {
			got := T(k)
			if want != got {
				t.Errorf("Compare Bundle with JSON:\ngot  %v\nwant %v", got, want)
			}
		}
	})

	// Test with ja.json
	fn := filepath.Join("testdata", "l10n.ja.json")
	testCases = map[string]string {
		"rootCmdShort": "rootCmdShort JaMessage",
		"rootCmdLong": "rootCmdLong JaMessage",
	}

	t.Run("T with ja.json", func(t *testing.T) {
		// Test Run
		err := LoadTranslation(fn)
		if err != nil {
			t.Fatalf("Failed to initialize: %s\n", err)
		}
		for k, want := range testCases {
			got := T(k)
			if want != got {
				t.Errorf("Compare Bundle with JSON:\ngot  %v\nwant %v", got, want)
			}
		}
	})
}

// Test for func NormalizeLang(
func TestNormalizeLang(t *testing.T) {
	testCases := []struct {
		lang string
		a string
		b string
	}{
		{
			lang: "ja-JP",
			a: "ja",
			b: "JP",
		},
		{
			lang: "ja_JP.UTF-8",
			a: "ja",
			b: "JP",
		},
	}
	t.Run("NormalizeLang", func(t *testing.T) {
		for _, want := range testCases {
			gotA, gotB := NormalizeLang(want.lang)
			if gotA != want.a {
				t.Errorf("Compare:\ngot  %v\nwant %v", gotA, want.a)
			}
			if gotB != want.b {
				t.Errorf("Compare:\ngot  %v\nwant %v", gotB, want.b)
			}
		}
	})
}
