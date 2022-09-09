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
	"hugo-docs-i18n/doci18n"
)

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
		datadir := viper.GetString("data-dir")
		infn := viper.GetString("md")
		if infn == "" {
			infn = viper.GetString("iso-md")
		}
		outfn := viper.GetString("json")
		if outfn == "" {
			outfn = viper.GetString("iso-json")
		}
		var infname, outfname string
		if !doci18n.IsExist(infn) {
			infname = filepath.Join(datadir, infn)
			outfname = filepath.Join(datadir, outfn)
			if !doci18n.IsExist(infname) {
				infname = ""
			}
		} else {
			infname = infn
			outfname = outfn
		}
		if infname == "" {
			log.Fatalf(`Locale JSON file dose not exist: %s`, infn)
		}
		// convert Markdown to JSON
		fmt.Printf("converting %s to %s ...\n", infname, outfname)
		if err := locale.ConvertLocaleFile(infname, outfname); err != nil {
			log.Fatalf("Error when converting: %w\n", err)
		}
		/* Save config file
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Error when writing config file: %w\n", err)
		}
		*/
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command and all subcommands,
	// localedbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called.
	convertCmd.Flags().StringP("md", "m", "", "source Markdown file of locale database")
	convertCmd.Flags().StringP("json", "j", "", "locale database file that created")
	convertCmd.Flags().Lookup("md").NoOptDefVal = ""
	convertCmd.Flags().Lookup("json").NoOptDefVal = ""
	viper.BindPFlag("md", convertCmd.Flags().Lookup("md"))
	viper.BindPFlag("json", convertCmd.Flags().Lookup("json"))
}
