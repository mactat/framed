/*
Copyright © 2023 Maciej Tatarski maciektatarski@gmail.com
*/

// Package cmd represents the command line interface of the application
package cmd

import (
	"fmt"
	"framed/pkg/ext"
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
		_, dirList := ext.ReadConfig(path)

		allGood := true
		// verify directories
		for _, dir := range dirList {
			if !ext.DirExists(dir.Path) {
				ext.PrintOut("❌ Directory not found ==>", dir.Path)
				allGood = false
			}

			// verify files
			if dir.Files != nil {
				verifyFiles(dir, &allGood)
			}

			// verify minCount and maxCount
			numFiles := ext.CountFiles(dir.Path)
			if numFiles < dir.MinCount {
				ext.PrintOut("❌ Min count ("+fmt.Sprint(dir.MinCount)+") not met ==>", dir.Path)
				allGood = false
			}
			if numFiles > dir.MaxCount {
				ext.PrintOut("❌ Max count ("+fmt.Sprint(dir.MaxCount)+") exceeded ==>", dir.Path)
				allGood = false
			}

			// verify childrenAllowed
			if !dir.AllowChildren {
				if ext.HasDirs(dir.Path) {
					ext.PrintOut("❌ Children not allowed ==>", dir.Path)
					allGood = false
				}
			}

			// verify maxDepth
			if ext.CheckDepth(dir.Path) > dir.MaxDepth {
				ext.PrintOut("❌ Max depth exceeded ("+fmt.Sprint(dir.MaxDepth)+") ==>", dir.Path)
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

func verifyFiles(dir ext.SingleDir, allGood *bool) {
	for _, file := range *dir.Files {
		if !ext.FileExists(dir.Path + "/" + file) {
			ext.PrintOut("❌ File not found      ==>", dir.Path+"/"+file)
			*allGood = false
		}
	}
}
func verifyForbiddenPatterns(dir ext.SingleDir, allGood *bool) {
	for _, pattern := range *dir.ForbiddenPatterns {
		matched := ext.MatchPatternInDir(dir.Path, pattern)
		for _, match := range matched {
			ext.PrintOut("❌ Forbidden pattern ("+pattern+") matched under ==>", dir.Path+"/"+match)
			*allGood = false
		}
	}
}

func verifyAllowedPatterns(dir ext.SingleDir, allGood *bool) {
	matchedCount := 0
	for _, pattern := range *dir.AllowedPatterns {
		matched := ext.MatchPatternInDir(dir.Path, pattern)
		matchedCount += len(matched)
	}
	if matchedCount != ext.CountFiles(dir.Path) && len(*dir.AllowedPatterns) > 0 {
		patternsString := strings.Join(*dir.AllowedPatterns, " ")
		ext.PrintOut("❌ Not all files match required pattern ("+patternsString+") under ==>", dir.Path)
		*allGood = false
	}
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.
	testCmd.PersistentFlags().String("template", "./framed.yaml", "path to template file")
}
