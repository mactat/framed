/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
			err := importFromUrl(output, url)
			if err != nil {
				fmt.Println("Error importing from url: ", err)
				return
			}
		}

		if example != "" {
			err := importFromUrl(output, exampleToUrl(example))
			if err != nil {
				fmt.Println("Error importing from example: ", err)
				return
			}
		}

		print("✅ Saved to ==>", output)

		// try to load
		configTree, _ := readConfig(output)
		print("✅ Imported successfully ==>", configTree.Name)

	},
}

func exampleToUrl(example string) string {
	return "https://raw.githubusercontent.com/mactat/framed/master/examples/" + example + ".yaml"
}

func importFromUrl(path string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
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
