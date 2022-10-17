/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Markdown locale
*/

package locale

import (
	// "bufio"
	"fmt"
	// "log"
	"os"
	"encoding/json"
	"regexp"
	// "unicode/utf8"
)

// Locale (lang/code pair) data
type LangPair struct {
	Lang string `json:"lang"`
	Code string `json:"code"`
}

type LangJSON struct {
	Locale []LangPair `json:"locale"`
}

// Search Locale JSON file
func SearchLocaleFile(infn, lg, cd string) ([]LangPair, error) {
	var pairs []LangPair
	// read locale (lang/code pair) data from JSON file
	contentBytes, err := os.ReadFile(infn)
	if err != nil {
		// return nil, fmt.Errorf("Error when reading file: %s", err)
		return nil, err
	}
	var gotJson LangJSON
	if err := json.Unmarshal(contentBytes, &gotJson); err != nil {
		// return nil, fmt.Errorf("Error when unmarshaling: %s", err)
		return nil, err
	}

	// Match flag?
	for _, lp := range gotJson.Locale {
		var matchLang, matchCode bool
		if lg != "" {
			matchLang, _ = regexp.MatchString(lg, lp.Lang)
		}
		if cd != "" {
			matchCode, _ = regexp.MatchString(cd, lp.Code)
		}
		if matchLang || matchCode || (lg == "" && cd == "") {
			pairs = append(pairs, LangPair{Lang: lp.Lang, Code: lp.Code})
			// fmt.Printf("Language: %30s \tCode: %5s\n", lp.Lang, lp.Code)
		}
	}
	return pairs, nil
}

// Print searched result
func PrintSearchedResult(pairs []LangPair)  {
	for _, lp := range pairs {
		fmt.Printf("Language: %30s \tCode: %5s\n", lp.Lang, lp.Code)
	}
}
