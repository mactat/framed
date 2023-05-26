/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func createDir(path string) {
	// Create directory
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		fmt.Printf("%-35s %-35s\n", "📂 Creating directory ==> ", path)
		if err != nil {
			log.Println(err)
		}
	}

}

func createFile(path string, name string) {
	// Check if file exists
	if _, err := os.Stat(path + "/" + name); errors.Is(err, os.ErrNotExist) {
		// Create file
		fmt.Printf("%-35s %-35s\n", "💽 Creating file      ==> ", path+"/"+name)
		file, err := os.Create(path + "/" + name)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
	}
}

// Check if directory exists on given path and is type dir
func dirExists(path string) bool {
	if path == "." {
		return true
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Check if file exists on given path and is type file
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Count files in given directory
func countFiles(path string) int {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return len(files)
}

func hasDirs(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			return true
		}
	}
	return false
}

// Check depth of folder tree, exclude .git folder
func checkDepth(path string) int {
	maxDepth := 0
	var depth int
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		} else if info.IsDir() {
			depth = strings.Count(path, string(os.PathSeparator)) + 1
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return maxDepth
}
