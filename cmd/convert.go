/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	// "os"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"hugo-docs-i18n/locale"
)

// Const & Variables
const dir = "data"
const infn = "ISO_639-1.md"
const outfn = "ISO_639-1.json"
var inFilename string
var outFilename string
var cfgF string

// localedbCmd represents the localedb command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Create Locale Database from ISO 639-1 standard language codes",
	Long: `hugo-docs-i18n init command need locale code for translated language.
This command create locale database from the Markdown file of the follows:

ISO 639-1 standard language codes:
https://www.andiamo.co.uk/resources/iso-language-codes/

The above file is saved as data/ISO_639-1.md,
and will be converted to JSON file as data/ISO_639-1.json.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("convert called.")
		if err := locale.ConvertLocaleFile(inFilename, outFilename); err != nil {
			log.Fatalf("Error when converting file: %s\n  %s -> %s", err, inFilename, outFilename)
		}
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Error when writing config file: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.
	infname := filepath.Join(dir, infn)
	outfname := filepath.Join(dir, outfn)

	// Cobra supports Persistent Flags which will work for this command and all subcommands,
	// localedbCmd.PersistentFlags().String("foo", "", "A help for foo")
	convertCmd.PersistentFlags().StringVar(&outFilename, "locale-db", outfname, "locale database file that created")
	viper.BindPFlag("locale-db", convertCmd.PersistentFlags().Lookup("locale-db"))

	// Cobra supports local flags which will only run when this command is called.
	convertCmd.Flags().StringVar(&inFilename, "src-md", infname, "source Markdown file of locale database")
	viper.BindPFlag("src-md", convertCmd.Flags().Lookup("src-md"))
}
