/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

// Package cmd represents the command line interface of the application
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
	Short: "Verify the project rules and structure",
	Long: `This command is verifying the project structure for consistency and compliance with the YAML template.
	
Example:
framed verify --template ./framed.yaml
`,
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("template").Value.String()
		// read config
		_, dirList := readConfig(path)

		allGood := true
		// verify directories
		for _, dir := range dirList {
			if !dirExists(dir.Path) {
				print("❌ Directory not found ==>", dir.Path)
				allGood = false
			}

			// verify files
			if dir.Files != nil {
				verifyFiles(dir, &allGood)
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
			if dir.ForbiddenPatterns != nil {
				verifyForbiddenPatterns(dir, &allGood)
			}

			// Verify allowed patterns
			if dir.AllowedPatterns != nil {
				verifyAllowedPatterns(dir, &allGood)
			}
		}

		if allGood {
			fmt.Println("✅ Verified successfully!")
		} else {
			os.Exit(1)
		}
	},
}

func verifyFiles(dir SingleDir, allGood *bool) {
	for _, file := range *dir.Files {
		if !fileExists(dir.Path + "/" + file) {
			print("❌ File not found      ==>", dir.Path+"/"+file)
			*allGood = false
		}
	}
}
func verifyForbiddenPatterns(dir SingleDir, allGood *bool) {
	for _, pattern := range *dir.ForbiddenPatterns {
		matched := matchPatternInDir(dir.Path, pattern)
		for _, match := range matched {
			print("❌ Forbidden pattern ("+pattern+") matched under ==>", dir.Path+"/"+match)
			*allGood = false
		}
	}
}

func verifyAllowedPatterns(dir SingleDir, allGood *bool) {
	matchedCount := 0
	for _, pattern := range *dir.AllowedPatterns {
		matched := matchPatternInDir(dir.Path, pattern)
		matchedCount += len(matched)
	}
	if matchedCount != countFiles(dir.Path) && len(*dir.AllowedPatterns) > 0 {
		patternsString := strings.Join(*dir.AllowedPatterns, " ")
		print("❌ Not all files match required pattern ("+patternsString+") under ==>", dir.Path)
		*allGood = false
	}
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.
	testCmd.PersistentFlags().String("template", "./framed.yaml", "path to template file")
}
