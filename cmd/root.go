/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}
