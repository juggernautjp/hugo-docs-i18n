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
	"github.com/spf13/viper"
)

// Variables
var dstLocale string
var contentDir string
var languageName string
var weight int
var time_format_default string
var time_format_blog string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Hugo ドキュメント翻訳環境の初期化",
	Long: `Initialize the configuration file for localization of Hugo Documentation.
For example of Japanese localization:

$ hugo-docs-i18n init --locale ja`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called: ", dstLocale)
		if dstLocale == "" {
			log.Fatalln(`--locale flag should be specified.
	"search" command can search locale code for the translation language. `)
		}
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Error when writing config file: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Persistent Flags which will work for this command and all subcommands.
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	initCmd.PersistentFlags().StringVar(&dstLocale, "locale", "", "locale code for target language")
	initCmd.MarkFlagRequired("locale")
	viper.BindPFlag("locale", initCmd.PersistentFlags().Lookup("locale"))

	// local flags which will only run when this command is called directly.
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().StringVar(&contentDir, "content-dir", "content", "content directory for target language")
	initCmd.Flags().StringVar(&languageName, "lang-name", "", "name of target language")
	initCmd.Flags().IntVar(&weight, "weight", 2, "weight")
	initCmd.Flags().StringVar(&time_format_default, "time-format-default", "2006/02/01", "default time format")
	initCmd.Flags().StringVar(&time_format_blog, "time-format-blog", "2006/02/01", "blog's time format")
	viper.BindPFlag("content-dir", initCmd.Flags().Lookup("content-dir"))
	viper.BindPFlag("lang-name", initCmd.Flags().Lookup("lang-name"))
	viper.BindPFlag("weight", initCmd.Flags().Lookup("weight"))
	viper.BindPFlag("time-format-default", initCmd.Flags().Lookup("time-format-default"))
	viper.BindPFlag("time-format-blog", initCmd.Flags().Lookup("time-format-blog"))
}
