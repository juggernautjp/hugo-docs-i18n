/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Markdown locale
*/

package locale

import (
	"bufio"
	"fmt"
	// "log"
	"os"
	"regexp"
	// "unicode/utf8"
)

// Constant
const regExp = `^\|\s+(?P<lang>(?:[-a-zA-Z]+)(?:\s\((?:[a-zA-Z\.]+)(?:\s[a-zA-Z]+)*\))?)\s+\|\s+(?P<code>(?:[a-z]{2})(?:-[a-z]{2})?)\s+\|$`
const template = "  { \"lang\": \"$lang\", \"code\": \"$code\" },\n"
const header = "{ \"locale\": [\n"
const footer = "{}\n]}\n"


// Convert Locale Markdown file to JSON file
func ConvertLocaleFile(infn, outfn string) error {
	// open the file
	infile, err := os.Open(infn)
	if err != nil {
		// return fmt.Errorf("Error when opening file: %w", err)
		return err
	}
	// close the file
	defer infile.Close()

	// open the file
	outfile, err := os.Create(outfn)
	if err != nil {
		// return fmt.Errorf("Error when opening file: %w", err)
		return err
	}
	// close the file
	defer outfile.Close()
	// write header to file
	fmt.Fprint(outfile, header)

	// read line by line
	s := bufio.NewScanner(infile)
	for s.Scan() {
		var line string = s.Text()
		// fmt.Println(line)
		// Regex pattern captures "key: value" pair from the content.
		pattern := regexp.MustCompile(regExp)

		result := []byte{}
		// For each match of the regex in the content.
		for _, submatches := range pattern.FindAllStringSubmatchIndex(line, -1) {
			// Apply the captured submatches to the template and append the output to the result.
			result = pattern.ExpandString(result, template, line, submatches)
		}
		// fmt.Println(string(result))
		fmt.Fprint(outfile, string(result))
	}
	if err := s.Err(); err != nil {
		// return fmt.Errorf("Error while reading file: %w", err)
		return err
	}
	// write footer to file
	fmt.Fprint(outfile, footer)
	return nil
}
