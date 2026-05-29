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

		report, err := ext.Verify(path)
		if err != nil {
			os.Exit(1)
		}

		for _, f := range report.Failures {
			ext.PrintOut(f.Prompt, f.Detail)
		}

		if report.OK {
			fmt.Println("✅ Verified successfully!")
			return
		}
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.PersistentFlags().String("template", "./framed.yaml", "path to template file")
}
