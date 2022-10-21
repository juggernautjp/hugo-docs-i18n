/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
	"github.com/juggernautjp/hugo-docs-i18n/doci18n"
)

// semverCmd represents the semver command
var semverCmd = &cobra.Command{
	Use:   "semver",
	Short: "Set version with SemVer format",
	Long: `Set the version of hugo-docs-i18n with SemVer format of "<Major>.<Minor>.<Patch>".`,
	Run: func(cmd *cobra.Command, args []string) {
		if !doci18n.CanSetVersionInfo() {
			log.Fatalf("Error: can not run semver subcommand.\nsemver subcommand is available for developers.\n")
		}
		doci18n.LoadVersionInfo()
		// inc flag
		kindi, _ := cmd.Flags().GetString("inc")
		switch kindi {
		case "ma", "major":
			doci18n.IncMajor()
		case "mi", "minor":
			doci18n.IncMinor()
		case "pa", "patch":
			doci18n.IncPatch()
		}
		// prerelease flag
		pre, _ := cmd.Flags().GetString("prerelease")
		if pre != "" {
			doci18n.SetPrerelease(pre)
		}
		// metabuild flag
		meta, _ := cmd.Flags().GetString("metabuild")
		if meta != "" {
			doci18n.SetMetabuild(meta)
		}
		if kindi != "" || pre != "" || meta != "" {
			fmt.Printf("New version: %s\n", doci18n.GetSemver())
			if err := doci18n.SaveVersionInfo(); err != nil {
				log.Fatalf("can not save version.json: %s\n", err)
			}
			os.Exit(0)
		}
		// force flag
		ver, _ := cmd.Flags().GetString("force")
		if ver == "" {
			log.Fatalf("--force is not specified: %s\n", ver)
		}
		err := doci18n.SetSemver(ver)
		if err != nil {
			log.Fatalf("--force is not SemVer format: %s\n", ver)
		} else {
			fmt.Printf("New version: %s\n", doci18n.GetSemver())
			if err := doci18n.SaveVersionInfo(); err != nil {
				log.Fatalf("can not save version.json: %s\n", err)
			}
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(semverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// semverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// semverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	semverCmd.Flags().StringP("inc", "i", "", "Increment version with SemVer format")
	semverCmd.Flags().Lookup("inc").NoOptDefVal = ""
	semverCmd.Flags().StringP("prerelease", "p", "", "Set Prerelease of SemVer")
	semverCmd.Flags().Lookup("prerelease").NoOptDefVal = ""
	semverCmd.Flags().StringP("metabuild", "m", "", "Set Metabuild of SemVer")
	semverCmd.Flags().Lookup("metabuild").NoOptDefVal = ""
	semverCmd.Flags().StringP("force", "f", "", "Force to set version with SemVer format")
	semverCmd.Flags().Lookup("force").NoOptDefVal = ""
}
