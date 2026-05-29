/*
Copyright © 2023 Maciej Tatarski maciektatarski@gmail.com
*/

// Package cmd represents the command line interface of the application
package cmd

import (
	"framed/pkg/ext"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project structure using a YAML template",
	Long: `This command is creating a new project structure from a YAML template.

Example:
framed create --template ./framed.yaml --files true
	`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := ext.CreateOptions{
			TemplatePath: cmd.Flag("template").Value.String(),
			CreateFiles:  cmd.Flag("files").Value.String() == "true",
		}
		if _, err := ext.Create(opts); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().String("template", "./framed.yaml", "path to template file default")

	createCmd.PersistentFlags().Bool("files", false, "create required files")
}
