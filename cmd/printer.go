/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "fmt"

func print(prompt string, text string) {
	fmt.Printf("%-35s %-35s\n", prompt, text)
}
