/*
Test getlang_*.go
*/

package locale

import (
	// "os"
	// "log"
	// "encoding/json"
	// "path/filepath"
	// "reflect"
	"testing"
)

// Test for InitTranslation()
func TestGetLang(t *testing.T) {
	var want string
	if IsBash() {
		want = "ja_JP.UTF-8"
	} else {
		want = "ja-JP"
	}
	t.Run("GetLang", func(t *testing.T) {
		// GetLang()
		got := GetLang()
		if got != want {
			t.Errorf("\ngot: %s\nwant: %s\n", got, want)
		}
	})
}
