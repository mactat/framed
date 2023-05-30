/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize the project structure",
	Long: `This command is visualizing the project structure from a YAML template.

Example:
framed visualize --template ./framed.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("template").Value.String()

		// read config
		_, dirList := readConfig(path)

		// visualize template
		visualizeTemplate(dirList)
	},
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	// Here you will define your flags and configuration settings.

	visualizeCmd.PersistentFlags().String("template", "./framed.yaml", "Path to template file default is ./framed.yaml")
}
