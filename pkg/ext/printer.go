package ext

import (
	"fmt"
	"os"
	"strings"

	"github.com/TwiN/go-color"
)

func PrintOut(prompt string, text string) {
	fmt.Printf("%-35s %-35s\n", prompt, text)
}

// This is ugly but it works, it needs to be refactored. It also has some bugs in case of out of order directories.
func VisualizeTemplate(template []SingleDir) {
	for dirNum, dir := range template {
		connectorDir := "├──"
		initString := ""
		depth := strings.Count(dir.Path, string(os.PathSeparator)) + 1
		name := strings.Split(dir.Path, string(os.PathSeparator))[depth-1]

		dirDepth := depth
		if depth <= 2 {
			dirDepth = 1
		} else {
			initString = "│"
		}

		if dirNum == len(template)-1 {
			connectorDir = "└──"
		}

		printDirectory(initString, dirDepth, connectorDir, name)

		if dir.Files == nil {
			continue
		}

		for num, file := range *dir.Files {
			connector := "├──"
			if depth > 1 {
				initString = "│"
			}

			if num == len(*dir.Files)-1 {
				connector = "└──"
			}

			printFile(initString, depth, connector, file)
		}
	}
}

func printDirectory(initString string, dirDepth int, connectorDir string, name string) {
	output := initString + strings.Repeat("    ", dirDepth-1) + connectorDir + " 📂 " + color.Ize(color.Blue, name)
	println(output)
}

func printFile(initString string, depth int, connector string, file string) {
	output := initString + strings.Repeat("    ", depth-1) + connector + " 📄 " + color.Ize(color.Green, file)
	println(output)
}
