/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

func exportConfig(name string, path string, subdirs []string, files []string, patterns map[string]string) {
	// create config, files and dirs are empty
	config := Config{
		Name: name,
		Structure: &SingleDir{
			Name:  "root",
			Files: &[]string{},
			Dirs:  &[]SingleDir{},
		},
	}
	// add subdirs
	insertSubdirs(config.Structure.Dirs, subdirs)

	// add files
	insertFiles(config.Structure, files)

	// add patterns
	insertPatterns(config.Structure, patterns)

	// export config
	yamlFile, err := yaml.Marshal(config)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	// Save to file
	err = ioutil.WriteFile(path, yamlFile, 0644)
	if err != nil {
		fmt.Printf("Error while writing file. %v", err)
	}
}

func insertSubdirs(dirs *[]SingleDir, subdirs []string) {
	for subdir := range subdirs {
		insertSingleDir(dirs, subdirs[subdir])
	}
}

func insertSingleDir(dirs *[]SingleDir, dir string) {
	subdirPath := strings.Split(dir, "/")
	curDirs := dirs
	for i := range subdirPath {
		if !containsDir(*curDirs, subdirPath[i]) {
			*curDirs = append(*curDirs, SingleDir{
				Name:            subdirPath[i],
				Files:           &[]string{},
				Dirs:            &[]SingleDir{},
				AllowedPatterns: &[]string{},
			})
		}
		// go deeper
		curDirs = getDir(*curDirs, subdirPath[i]).Dirs
	}
}

func containsDir(dirs []SingleDir, name string) bool {
	for _, dir := range dirs {
		if dir.Name == name {
			return true
		}
	}
	return false
}

func getDir(dirs []SingleDir, name string) *SingleDir {
	for _, dir := range dirs {
		if dir.Name == name {
			return &dir
		}
	}
	return nil
}

func insertFiles(root *SingleDir, files []string) {
	for file := range files {
		insertSingleFile(root, files[file])
	}
}

func insertSingleFile(root *SingleDir, file string) {
	subdirPath := strings.Split(file, "/")

	if len(subdirPath) == 1 {
		*root.Files = append(*root.Files, subdirPath[0])
		return
	}
	curDirs := root.Dirs
	curDir := root
	for i := 0; i < len(subdirPath)-1; i++ {
		curDir = getDir(*curDirs, subdirPath[i])
		curDirs = curDir.Dirs
	}
	*curDir.Files = append(*curDir.Files, subdirPath[len(subdirPath)-1])

}

func insertPatterns(root *SingleDir, patterns map[string]string) {
	for dir, pattern := range patterns {
		insertSinglePattern(root, dir, pattern)
	}
}

func insertSinglePattern(root *SingleDir, dir string, pattern string) {
	subdirPath := strings.Split(dir, "/")
	curDirs := root.Dirs
	curDir := root
	for i := range subdirPath {
		// go deeper
		curDir = getDir(*curDirs, subdirPath[i])
		curDirs = curDir.Dirs
	}
	// insert pattern
	*curDir.AllowedPatterns = append(*curDir.AllowedPatterns, pattern)
}
