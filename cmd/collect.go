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
	// "os"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
	"github.com/juggernautjp/hugo-docs-i18n/doci18n"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect statistics of translation progress",
	Long: `Output the number of translated/not-translated files, 
and what percentage of the files have been translated.

The data is stored "data/i18n/<code>.json, and show the page.`,
	Run: func(cmd *cobra.Command, args []string) {
		// tdir := viper.GetString("target-dir")
		// tcode := viper.GetString("target-code")
		tdir, _ := cmd.Flags().GetString("target-dir")
		tcode, _ := cmd.Flags().GetString("target-code")
		if tdir == "" && tcode == "" {
			log.Fatalln(`Both flags are not specified`)
		}
		if tcode == "" {
			tcode = filepath.Base(tdir)
		} else if tdir == "" {
			// cdir := viper.GetString("content-dir")
			cdir, _ := cmd.Flags().GetString("content-dir")
			tdir = filepath.Join(cdir, tcode)
		}
		fmt.Printf("count the number of draft files under the directory: %s\n", tdir)
		// Get target directory name
		/*
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf(`Can not get current directory: %w`, err)
		}
		targetdname := filepath.Join(pwd, tdir)
		*/
		if !doci18n.IsExist(tdir) {
			log.Fatalf(`Directory dose not exist: %s`, tdir)
		}
		// Count draft files under the directory
		gotPJ, gotCD, err := doci18n.CountDraftFile(tdir)
		if err != nil {
			log.Fatalf(`Can not count draft files: %s`, err)
		}
		// Get data filename to save JSON data
		// ddir := viper.GetString("data-dir")
		ddir, _ := cmd.Flags().GetString("data-dir")
		jsonfn := filepath.Join(ddir, tcode + ".json")
		if err := doci18n.SaveCountJSONFile(jsonfn, gotPJ, gotCD); err != nil {
			log.Fatalf("Can not save data: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	collectCmd.Flags().StringP("target-dir", "d", "", "Content directory for target language")
	collectCmd.Flags().StringP("target-code", "c", "", "Locale code for target language")
	collectCmd.Flags().Lookup("target-dir").NoOptDefVal = ""
	collectCmd.Flags().Lookup("target-code").NoOptDefVal = ""
	// viper.BindPFlag("target-dir", collectCmd.Flags().Lookup("target-dir"))
	// viper.BindPFlag("target-code", collectCmd.Flags().Lookup("target-code"))
}
