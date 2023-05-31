/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

// Package cmd represents the command line interface of the application
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
		fmt.Printf("%-35s %-35s\n", "📄 Creating file      ==> ", path+"/"+name)
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
	filesCount := 0
	for _, file := range files {
		if !file.IsDir() {
			filesCount++
		}
	}
	return filesCount
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

func matchPatternInDir(path string, pattern string) []string {
	if pattern == "" {
		pattern = ".*"
	}
	// List all files in directory
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	matched := []string{}
	for _, file := range files {
		if !file.IsDir() {
			match, err := regexp.MatchString(pattern, file.Name())
			if err != nil {
				log.Fatal(err)
			}
			if match {
				matched = append(matched, file.Name())
			}
		}
	}
	return matched
}

// Capture all subdirectories in given directory
func captureSubDirs(path string) []string {
	var dirs []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		} else if info.IsDir() && info.Name() != "." {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs
}

// Capture all files in given directory
func captureAllFiles(path string) []string {
	var files []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		} else if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// Capture rules for files with same extension in given directory. If all files in subdirectory have the same extension, save the extension to map with directory path as key.
// It should return map path -> extension
func captureRequiredPatterns(path string) map[string]string {
	var rules = make(map[string]string)
	var dirs []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		} else if info.IsDir() && info.Name() != "." {
			dirs = append(dirs, path)
		}
		return nil
	})
	// Check files in dir, if all extensions are the same, save extension to map with dir path as key
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		ext := ""
		for _, file := range files {
			if !file.IsDir() {
				extension := filepath.Ext(file.Name())
				if ext == "" {
					ext = extension
				} else if ext != extension {
					ext = ""
					break
				}
			}
		}
		if ext != "" {
			rules[dir] = ext
		}
	}
	return rules
}
