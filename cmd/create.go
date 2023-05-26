/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project structure using a YAML template",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("template").Value.String()
		createFiles := cmd.Flag("files").Value.String() == "true"

		// read config
		dirsList := readConfig(path)

		// create directories
		for _, dir := range dirsList {
			createDir(dir.Path)
		}

		// create files
		if createFiles {
			for _, dir := range dirsList {
				for _, file := range dir.Required.Files {
					createFile(dir.Path, file)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	createCmd.PersistentFlags().String("template", "./framed.yaml", "Path to template file default is ./framed.yaml")
	// add flag to create required files
	createCmd.PersistentFlags().Bool("files", false, "Create required files")
}
