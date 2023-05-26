/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify the project structure for consistency and compliance with the YAML template",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("template").Value.String()
		// read config
		dirsList := readConfig(path)

		allDirsFound := true
		// verify directories
		for _, dir := range dirsList {
			if !dirExists(dir.Path) {
				print("❌ Directory not found ==>", dir.Path)
				allDirsFound = false
			}
		}

		allFilesFound := true
		// verify files
		for _, dir := range dirsList {
			for _, file := range dir.Required.Files {
				if !fileExists(dir.Path + "/" + file) {
					print("❌ File not found      ==>", dir.Path+"/"+file)
					allFilesFound = false
				}
			}
		}

		if allDirsFound && allFilesFound {
			println("✅ All directories and files found")
		} else {
			println("\n❌ Some directories or files not found")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.
	testCmd.PersistentFlags().String("template", "./framed.yaml", "Path to template file default is ./framed.yaml")
}
