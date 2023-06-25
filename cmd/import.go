/*
Copyright © 2023 Maciej Tatarski maciektatarski@gmail.com
*/

// Package cmd represents the command line interface of the application
package cmd

import (
	"fmt"
	"framed/pkg/ext"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import the project structure",
	Long: `This command is importing the project structure from a YAML template. It can be imported from a template or from a remote URL.
Example:
framed import https://raw.githubusercontent.com/username/repo/master/framed.yaml
or
framed import --example python --output ./python.yaml
`,
	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		example := cmd.Flag("example").Value.String()
		output := cmd.Flag("output").Value.String()

		if url != "" {
			err := ext.ImportFromUrl(output, url)
			if err != nil {
				fmt.Println("Error importing from url: ", err)
				return
			}
		}

		if example != "" {
			err := ext.ImportFromUrl(output, ext.ExampleToUrl(example))
			if err != nil {
				fmt.Println("Error importing from example: ", err)
				return
			}
		}

		ext.PrintOut("✅ Saved to ==>", output)

		// try to load
		configTree, _ := ext.ReadConfig(output)
		ext.PrintOut("✅ Imported successfully ==>", configTree.Name)

	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.
	importCmd.PersistentFlags().String("url", "", "url to template file")
	importCmd.PersistentFlags().String("example", "", "example template file from github")
	importCmd.MarkFlagsMutuallyExclusive("url", "example")

	// path to file
	importCmd.PersistentFlags().String("output", "./framed.yaml", "path to template file")
}
