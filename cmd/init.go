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

	"github.com/spf13/cobra"
	"github.com/juggernautjp/hugo-docs-i18n/doci18n"
)

// Const
const config_name = "hugo-docs-i18n.yaml"


// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Hugo ドキュメント翻訳環境の設定ファイルの初期化",
	Long: `Generate a configuration file for localization of Hugo Documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Generating config file: %s ...", config_name)
		// Generate initial config file
		if err := doci18n.SaveConfigFile(config_name); err != nil {
			log.Fatalf(`Can not generate config file: %s`, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Persistent Flags which will work for this command and all subcommands.
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	// initCmd.PersistentFlags().String("locale", "", "locale code for target language")
	// initCmd.MarkFlagRequired("locale")
	// viper.BindPFlag("locale", initCmd.PersistentFlags().Lookup("locale"))

	// local flags which will only run when this command is called directly.
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
