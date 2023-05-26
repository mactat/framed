/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

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

		allGood := true
		// verify directories
		for _, dir := range dirsList {
			if !dirExists(dir.Path) {
				print("❌ Directory not found ==>", dir.Path)
				allGood = false
			}

			// verify files
			for _, file := range dir.Required.Files {
				if !fileExists(dir.Path + "/" + file) {
					print("❌ File not found      ==>", dir.Path+"/"+file)
					allGood = false
				}
			}

			// verify minCount and maxCount
			numFiles := countFiles(dir.Path)
			if numFiles < dir.MinCount {
				print("❌ Min count ("+fmt.Sprint(dir.MinCount)+") not met ==>", dir.Path)
				allGood = false
			}
			if numFiles > dir.MaxCount {
				print("❌ Max count ("+fmt.Sprint(dir.MaxCount)+") exceeded ==>", dir.Path)
				allGood = false
			}

			// verify childrenAllowed
			if !dir.AllowChildren {
				if hasDirs(dir.Path) {
					print("❌ Children not allowed ==>", dir.Path)
					allGood = false
				}
			}

			// verify maxDepth
			if checkDepth(dir.Path) > dir.MaxDepth {
				print("❌ Max depth exceeded ("+fmt.Sprint(dir.MaxDepth)+") ==>", dir.Path)
				allGood = false
			}

			// Verify forbidden
			for _, pattern := range dir.Forbidden.Patterns {
				matched := matchPatternInDir(dir.Path, pattern)
				for _, match := range matched {
					print("❌ Forbidden pattern ("+pattern+") matched under ==>", dir.Path+"/"+match)
					allGood = false
				}
			}

			// Verify required
			matchedCount := 0
			for _, pattern := range dir.Required.Patterns {
				matched := matchPatternInDir(dir.Path, pattern)
				matchedCount += len(matched)
			}
			if matchedCount != countFiles(dir.Path) && len(dir.Required.Patterns) > 0 {
				patternsString := strings.Join(dir.Required.Patterns, " ")
				print("❌ Not all files match required pattern ("+patternsString+") under ==>", dir.Path)
				allGood = false
			}
		}

		if allGood {
			fmt.Println("✅ Verified successfully!")
		} else {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.
	testCmd.PersistentFlags().String("template", "./framed.yaml", "Path to template file default is ./framed.yaml")
}
