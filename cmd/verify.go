/*
Copyright © 2023 Maciej Tatarski maciektatarski@gmail.com
*/

// Package cmd represents the command line interface of the application
package cmd

import (
	"fmt"
	"framed/pkg/ext"
	"os"

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
				ext.VerifyFiles(dir, &allGood)
			}

			// verify minCount
			numFiles := ext.CountFiles(dir.Path)
			if numFiles < dir.MinCount {
				ext.PrintOut("❌ Min count ("+fmt.Sprint(dir.MinCount)+") not met ==>", dir.Path)
				allGood = false
			}

			// verify maxCount
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
				ext.VerifyForbiddenPatterns(dir, &allGood)
			}

			// Verify allowed patterns
			if dir.AllowedPatterns != nil {
				ext.VerifyAllowedPatterns(dir, &allGood)
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

	testCmd.PersistentFlags().String("template", "./framed.yaml", "path to template file")
}
