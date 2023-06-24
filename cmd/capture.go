/*
Copyright © 2023 Maciej Tatarski maciektatarski@gmail.com
*/

// Package cmd represents the command line interface of the application
package cmd

import (
	"fmt"

	"framed/pkg/ext"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture the current project structure as a YAML template",
	Long: `This command is capturing the current project structure as a YAML template.

Example:
framed capture --output ./framed.yaml --name my-project
`,
	Run: func(cmd *cobra.Command, args []string) {
		output := cmd.Flag("output").Value.String()
		name := cmd.Flag("name").Value.String()
		depthStr := cmd.Flag("depth").Value.String()
		depth, err := strconv.Atoi(depthStr)
		if err != nil {
			ext.PrintOut("🚨 Invalid depth value: ", depthStr)
			os.Exit(1)
		}
		ext.PrintOut("📝 Name:", name+"\n")

		// capture subdirectories
		subdirs := ext.CaptureSubDirs(".", depth)
		ext.PrintOut("📂 Directories:", fmt.Sprintf("%v", len(subdirs)))

		// capture files
		files := ext.CaptureAllFiles(".", depth)
		ext.PrintOut("📄 Files:", fmt.Sprintf("%v", len(files)))

		// capture patterns
		patterns := ext.CaptureRequiredPatterns(".", depth)
		ext.PrintOut("🔁 Patterns:", fmt.Sprintf("%v", len(patterns)))

		// export config
		ext.ExportConfig(name, output, subdirs, files, patterns)
		ext.PrintOut("\n✅ Exported to file: ", output)
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)

	// Here you will define your flags and configuration settings.
	captureCmd.PersistentFlags().String("output", "./framed.yaml", "path to output file")

	captureCmd.PersistentFlags().String("name", "default", "name of the project")

	// Int flag - depth
	captureCmd.PersistentFlags().String("depth", "-1", "depth of the directory tree to capture")
}
