/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/juggernautjp/hugo-docs-i18n/doci18n"
)

// semverCmd represents the semver command
var semverCmd = &cobra.Command{
	Use:   "semver",
	Short: "Set version with SemVer format",
	Long: `Set the version of hugo-docs-i18n with SemVer format of "<Major>.<Minor>.<Patch>".`,
	Run: func(cmd *cobra.Command, args []string) {
		svStr := viper.GetString("semver")
		if svStr == "" {
			log.Fatalf("You should execute \"hugo-docs-i18n init\"\n")
		}
		sv := doci18n.NewSemver()
		if err := sv.ForceValue(svStr); err != nil {
			log.Fatalf("semver in config file is not correct: %s\n", svStr)
		}
		// inc flag
		kindi, _ := cmd.Flags().GetString("inc")
		switch kindi {
		case "ma", "major":
			sv.IncMajor()
		case "mi", "minor":
			sv.IncMinor()
		case "pa", "patch":
			sv.IncMinor()
		}
		if kindi == "" {
			log.Fatalf("--inc is empty\n")
		} else {
			newv := sv.String()
			fmt.Printf("New version: %s\n", newv)
			viper.Set("semver", newv)
			viper.WriteConfig()
			os.Exit(0)
		}
		// prerelease flag
		pre, _ := cmd.Flags().GetString("prerelease")
		if pre == "" {
			log.Fatalf("--prerelease is empty\n")
		} else {
			sv.SetPrerelease(pre)
		}
		// metabuild flag
		meta, _ := cmd.Flags().GetString("metabuild")
		if meta == "" {
			log.Fatalf("--metabuild is empty\n")
		} else {
			sv.SetMetabuild(meta)
		}
		if pre != "" || meta != "" {
			newv := sv.String()
			fmt.Printf("New version: %s\n", newv)
			viper.Set("semver", newv)
			viper.WriteConfig()
			os.Exit(0)
		}
		// force flag
		ver, _ := cmd.Flags().GetString("force")
		if ver == "" {
			log.Fatalf("--force is empty\n")
		}
		err := sv.ForceValue(ver)
		if err != nil {
			log.Fatalf("--force is not SemVer format: %s\n", ver)
		} else {
			newv := sv.String()
			fmt.Printf("New version: %s\n", newv)
			viper.Set("semver", newv)
			viper.WriteConfig()
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
