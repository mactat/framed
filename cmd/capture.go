/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
		name := cmd.Flag("name").Value.String()
		print("📝 Name:", name+"\n")

		// capture subdirectories
		subdirs := captureSubDirs(".")
		print("📂 Directories:", fmt.Sprintf("%v", len(subdirs)))

		// capture files
		files := captureAllFiles(".")
		print("📄 Files:", fmt.Sprintf("%v", len(files)))

		// capture patterns
		patterns := captureRequiredPatterns(".")
		print("🔁 Patterns:", fmt.Sprintf("%v", len(patterns)))

		// export config
		exportConfig(name, output, subdirs, files, patterns)
		print("\n✅ Exported to file: ", output)
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)

	// Here you will define your flags and configuration settings.
	captureCmd.PersistentFlags().String("output", "./framed.yaml", "Path to output file")

	captureCmd.PersistentFlags().String("name", "default", "Name of the project")
}
