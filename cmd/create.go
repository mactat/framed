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
	Long: `This command is creating a new project structure from a YAML template.

Example:
framed create --template ./framed.yaml --files true
	`,
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("template").Value.String()
		createFiles := cmd.Flag("files").Value.String() == "true"

		// read config
		_, dirList := readConfig(path)

		// create directories
		for _, dir := range dirList {
			createDir(dir.Path)
		}

		// create files
		if createFiles {
			for _, dir := range dirList {
				if dir.Files == nil {
					continue
				}
				for _, file := range *dir.Files {
					createFile(dir.Path, file)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	createCmd.PersistentFlags().String("template", "./framed.yaml", "path to template file default")
	// add flag to create required files
	createCmd.PersistentFlags().Bool("files", false, "create required files")
}
