//go:build !windows

package locale

import (
	"os"
)

// Is Bash
func IsBash() bool {
	return true
}

// get Locale
func GetLang() string {
	return os.Getenv("LANG")
}
