//go:build windows

package locale

import (
	"io"
	"os"
	"fmt"
	"strings"

	"github.com/hnakamur/go-powershell"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// ShiftJIS から UTF-8
func sjis_to_utf8(str string) (string, error) {
	ret, err := io.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
} 


// Is Bash
func IsBash() bool {
	if os.Getenv("SHELL") != "" {
		return true
	} else {
		return false
	}
}

// get Locale
func GetLang() string {
	if IsBash() {
		// Bash
		return os.Getenv("LANG")
	} else {
		// PowerShell with go-powershell
		shell, err := powershell.New()
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		defer shell.Exit()
		out, err := shell.Exec("$(Get-Culture).Name")
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		utf8out, err := sjis_to_utf8(out)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		return strings.TrimRight(utf8out, "\r\n")
	}
}
