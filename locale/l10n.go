/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Markdown l10n
*/

package locale

import (
	"os"
	"fmt"
	// "bufio"
	// "log"
	"regexp"
	"path/filepath"
	"encoding/json"
	_ "embed"
)

// Constant
const regExp2 = `^([a-z]{2})[_-]([A-Z]{2})`
const dataDir = "data/i18n"
const tmpl1 = "l10n.%s.json"
const tmpl2 = "l10n.%s_%s.json"

//go:embed l10n.en.json
var contentBytes []byte
var EnMessages Bundle
var TargetMessages Bundle

// Bundle
type KeyPair struct {
		Key string `json:"key"`
		Message string `json:"message"`
}

type Bundle struct {
	Pairs []KeyPair `json:"translation"`
}

// find key/message pair
func (b Bundle) findKey(k string) string {
	for _, m := range b.Pairs {
		if m.Key == k {
			return m.Message
		}
	}
	return ""
}

func (b Bundle) isEmpty() bool {
	if len(b.Pairs) == 0 {
		return true
	} else {
		return false
	}
}


// File is exist?
func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// load translation
func LoadTranslation(fn string) error {
	if !isExist(fn) {
		return fmt.Errorf("File is not found: %s", fn)
	}

	// read translation data from JSON file
	cb, err := os.ReadFile(fn)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(cb, &TargetMessages); err != nil {
		return err
	}
	return nil
}

// load enMessages data from JSON
func InitTranslation() error {
	// fmt.Printf("Translation is initialized\n")
	if err := json.Unmarshal(contentBytes, &EnMessages); err != nil {
		return err
	}
	return nil
}

// Normalize LANG
func NormalizeLang(line string) (a, b string) {
	pattern := regexp.MustCompile(regExp2)
	var result []string
	// For each match of the regex in the content.
	result = pattern.FindStringSubmatch(line)
	return result[1], result[2]
}

// return translation
func T(k string) string {
	if EnMessages.isEmpty() {
		// initialize translation with en language
		InitTranslation()
	}
	// if translation exists
	line := GetLang()
	if line != "" {
		a, b := NormalizeLang(line)
		fn1 := fmt.Sprintf(tmpl1, a)
		path1 := filepath.Join(dataDir, fn1)
		fn2 := fmt.Sprintf(tmpl2, a, b)
		path2 := filepath.Join(dataDir, fn2)
		if isExist(path1) {
			LoadTranslation(path1)
		} else if isExist(path2) {
			LoadTranslation(path2)
		}
	}

	// Search message with key
	if TargetMessages.isEmpty() {
		return EnMessages.findKey(k)
	} 
	return TargetMessages.findKey(k)
}

