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

		opts := ext.CaptureOptions{
			Name:       name,
			OutputPath: output,
			Depth:      depth,
		}
		report, err := ext.CaptureWithReporter(opts, func(r ext.CaptureReport) {
			ext.PrintOut("📂 Directories:", fmt.Sprintf("%v", r.Dirs))
			ext.PrintOut("📄 Files:", fmt.Sprintf("%v", r.Files))
			ext.PrintOut("🔁 Patterns:", fmt.Sprintf("%v", r.Patterns))
		})
		if err != nil {
			os.Exit(1)
		}
		_ = report
		ext.PrintOut("\n✅ Exported to file: ", output)
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)

	captureCmd.PersistentFlags().String("output", "./framed.yaml", "path to output file")

	captureCmd.PersistentFlags().String("name", "default", "name of the project")

	captureCmd.PersistentFlags().String("depth", "-1", "depth of the directory tree to capture")
}
