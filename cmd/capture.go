/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture the current project structure as a YAML template",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		output := cmd.Flag("output").Value.String()
		fmt.Println(output)

		// capture subdirectories
		subdirs := captureSubDirs(".")
		fmt.Printf("Subdirs:\n\n%+v\n\n", strings.Join(subdirs, "\n"))

		// capture files
		files := captureAllFiles(".")
		fmt.Printf("Files:\n\n%+v\n\n", strings.Join(files, "\n"))

		// capture patterns
		patterns := captureRequiredPatterns(".")
		for dir, pattern := range patterns {
			fmt.Printf("Pattern Found: %s in %s\n", pattern, dir)
		}
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)

	// Here you will define your flags and configuration settings.
	captureCmd.PersistentFlags().String("output", "./framed.yaml", "Path to output file default is ./framed.yaml")
}
