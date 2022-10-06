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

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
	"github.com/juggernautjp/hugo-docs-i18n/doci18n"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Copy and update content files from English to target language",
	Long: `Copy and update content files written in English to content/<code> directory.
Before execute this command, you should run "hugo-docs-i18n init".`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
		// srcdir := viper.GetString("source-dir")
		srcdir, _ := cmd.Flags().GetString("source-dir")
		if srcdir == "" {
			log.Fatalln("You should execute \"hugo-docs-i18n init\"")
		}
		// dstdir := viper.GetString("target-dir")
		dstdir, _ := cmd.Flags().GetString("target-dir")
		if dstdir == "" {
			log.Fatalln("You should specify the target language directory with --target-dir flag.")
		}
		if !doci18n.IsExist(dstdir) {
			log.Fatalf(`Directory dose not exist: %s`, dstdir)
		}
		// Copy content directory from source language (English) to target language
		if _, err := doci18n.CopyContentDir(srcdir, dstdir); err != nil {
			log.Fatalf("Error when copying: %s -> %s: %s\n", srcdir, dstdir, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().StringP("target-dir", "d", "", "Content directory for target language")
	updateCmd.Flags().Lookup("target-dir").NoOptDefVal = ""
	// viper.BindPFlag("target-dir", updateCmd.Flags().Lookup("target-dir"))
}
