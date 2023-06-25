/*
Copyright © 2023 Maciej Tatarski maciektatarski@gmail.com
*/

// Package cmd represents the command line interface of the application
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "framed",
	Short: "CLI tool for managing folder and files structures",
	Long: `FRAMED (Files and Directories Reusability, Architecture, and Management)
is a powerful CLI tool written in Go that simplifies the organization and management
of files and directories in a reusable and architectural manner. It provides YAML
templates for defining project structures and ensures that your projects adhere to 
the defined structure, enabling consistency and reusability.

	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
