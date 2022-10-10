/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo Docs i18n Config file
*/

package doci18n

import (
	"os"
	"fmt"
)

// ExampleConfig is the config used within hugo-docs-i18n init.
const ExampleConfig = `# This is an example hugo-docs-i18n.yaml file with some sensible defaults.
# Make sure to check the documentation at https://github.com/juggernautjp/hugo-docs-i18n/
source-dir: content/en
content-dir: content
data-dir: data/i18n
iso-md: ISO_639-1.md
iso-json: ISO_639-1.json
semver: 0.1.0
ver-msg: Release version
# ver-msg: Dev version
`

// Save config data to YAML file
func SaveConfigFile(outfn string) error {
	// open the file
	conf, err := os.Create(outfn)
	if err != nil {
		return fmt.Errorf("Error when opening file: %w", err)
	}
	// close the file
	defer conf.Close()

	// Write `ExampleConfig` data to YAML config file 
	if _, err := conf.WriteString(ExampleConfig); err != nil {
		return fmt.Errorf("Error when writing file: %w", err)
	}
	return nil
}
