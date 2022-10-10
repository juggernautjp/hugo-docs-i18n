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
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/juggernautjp/hugo-docs-i18n/locale"
	"github.com/juggernautjp/hugo-docs-i18n/doci18n"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search locale code for target translation language",
	Long: `If you do not know your language locale code, you can search it with this command 
or the following web page:

ISO 639-1 standard language codes:
  https://www.andiamo.co.uk/resources/iso-language-codes/`,
	Run: func(cmd *cobra.Command, args []string) {
		// tl := viper.GetString("lang")
		// tc := viper.GetString("code")
		tl, _ := cmd.Flags().GetString("lang")
		tc, _ := cmd.Flags().GetString("code")
		fmt.Printf("searching lang=%s, code=%s ...\n", tl, tc)
		ldb := viper.GetString("iso-json")
		datadir := viper.GetString("data-dir")
		// ldb, _ := cmd.Flags().GetString("iso-json")
		// datadir, _ := cmd.Flags().GetString("data-dir")
		var fn string
		if !doci18n.IsExist(ldb) {
			fn = filepath.Join(datadir, ldb)
			if !doci18n.IsExist(fn) {
				fn = ""
			}
		} else {
			fn = ldb
		}
		if fn == "" {
			log.Fatalf(`Locale JSON file dose not exist: %s`, ldb)
		}
		pairs, err := locale.SearchLocaleFile(fn, tl, tc)
		if err != nil {
			log.Fatalf("Error when searching locale: %s\n", err)
		}
		locale.PrintSearchedResult(pairs)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	searchCmd.Flags().StringP("lang", "l", "", "name of target language")
	searchCmd.Flags().StringP("code", "c", "", "locale code of ISO 639-1")
	searchCmd.Flags().Lookup("lang").NoOptDefVal = ""
	searchCmd.Flags().Lookup("code").NoOptDefVal = ""
	// viper.BindPFlag("lang", searchCmd.Flags().Lookup("lang"))
	// viper.BindPFlag("code", searchCmd.Flags().Lookup("code"))
}
