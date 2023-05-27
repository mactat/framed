/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/TwiN/go-color"
)

func print(prompt string, text string) {
	fmt.Printf("%-35s %-35s\n", prompt, text)
}

// This is ugly but it works, it needs to be refactored. It also has some bugs in case of out of order directories.
func visualizeTemplate(template []SingleDir) {
	// Determine how deep is the path
	for dirNum, dir := range template {
		connectorDir := "├──"
		initString := ""
		depth := strings.Count(dir.Path, string(os.PathSeparator)) + 1
		name := strings.Split(dir.Path, string(os.PathSeparator))[depth-1]
		// Print the path with correct depth
		dirDepth := depth
		if depth <= 2 {
			dirDepth = 1
		}
		if depth > 2 {
			initString = "│"
		}
		if dirNum == len(template)-1 {
			connectorDir = "└──"
		}
		println(initString + strings.Repeat("    ", dirDepth-1) + connectorDir + " 📂 " + color.Ize(color.Blue, name))
		// Print the files
		for num, file := range dir.Required.Files {
			connector := "├──"
			if depth > 1 {
				initString = "│"
			}
			if num == len(dir.Required.Files)-1 {
				connector = "└──"
			}
			println(initString+strings.Repeat("    ", depth-1)+connector+" 💽", color.Ize(color.Green, file))
		}

	}
}
