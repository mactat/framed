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

// SingleDir struct
type SingleDirV2 struct {
	Name              string         `yaml:"name,omitempty"`
	Files             *[]string      `yaml:"files,omitempty"`
	Dirs              *[]SingleDirV2 `yaml:"dirs"`
	AllowedPatterns   *[]string      `yaml:"allowedPatterns,omitempty"`
	ForbiddenPatterns []string       `yaml:"forbiddenPatterns,omitempty"`
	MaxDepth          int            `yaml:"maxDepth,omitempty"`
	MaxCount          int            `yaml:"maxCount,omitempty"`
	MinCount          int            `yaml:"minCount,omitempty"`
}

// SingleDirV2 struct without Dir
type SingleDirV2WithoutDir struct {
	Name              string   `yaml:"name,omitempty"`
	Files             []string `yaml:"files,omitempty"`
	AllowedPatterns   []string `yaml:"allowedPatterns,omitempty"`
	ForbiddenPatterns []string `yaml:"forbiddenPatterns,omitempty"`
	MaxDepth          int      `yaml:"maxDepth,omitempty"`
	MaxCount          int      `yaml:"maxCount,omitempty"`
	MinCount          int      `yaml:"minCount,omitempty"`
}

// // Custom marshal function to ommits empty dirs
// func (d *SingleDirV2) MarshalYAML() (interface{}, error) {
// 	//fmt.Println(*d.Dirs)
// 	if d.Dirs == nil || len(*d.Dirs) == 0 {
// 		return SingleDirV2WithoutDir{
// 			Name:              d.Name,
// 			Files:             d.Files,
// 			AllowedPatterns:   d.AllowedPatterns,
// 			ForbiddenPatterns: d.ForbiddenPatterns,
// 			MaxDepth:          d.MaxDepth,
// 			MaxCount:          d.MaxCount,
// 			MinCount:          d.MinCount,
// 		}, nil
// 	}
// 	return &d, nil
// }

type Config struct {
	Name      string       `yaml:"name"`
	Structure *SingleDirV2 `yaml:"structure,omitempty"`
}

func exportConfig(name string, path string, subdirs []string, files []string, patterns map[string]string) {
	// create config, files and dirs are empty
	config := Config{
		Name: name,
		Structure: &SingleDirV2{
			Name:  "root",
			Files: &[]string{},
			Dirs:  &[]SingleDirV2{},
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

func insertSubdirs(dirs *[]SingleDirV2, subdirs []string) {
	for subdir := range subdirs {
		insertSingleDir(dirs, subdirs[subdir])
	}
}

func insertSingleDir(dirs *[]SingleDirV2, dir string) {
	subdirPath := strings.Split(dir, "/")
	curDirs := dirs
	for i := range subdirPath {
		if !containsDir(*curDirs, subdirPath[i]) {
			*curDirs = append(*curDirs, SingleDirV2{
				Name:            subdirPath[i],
				Files:           &[]string{},
				Dirs:            &[]SingleDirV2{},
				AllowedPatterns: &[]string{},
			})
		}
		// go deeper
		curDirs = getDir(*curDirs, subdirPath[i]).Dirs
	}
}

func containsDir(dirs []SingleDirV2, name string) bool {
	for _, dir := range dirs {
		if dir.Name == name {
			return true
		}
	}
	return false
}

func getDir(dirs []SingleDirV2, name string) *SingleDirV2 {
	for _, dir := range dirs {
		if dir.Name == name {
			return &dir
		}
	}
	return nil
}

func insertFiles(root *SingleDirV2, files []string) {
	for file := range files {
		insertSingleFile(root, files[file])
	}
}

func insertSingleFile(root *SingleDirV2, file string) {
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

func insertPatterns(root *SingleDirV2, patterns map[string]string) {
	for dir, pattern := range patterns {
		insertSinglePattern(root, dir, pattern)
	}
}

func insertSinglePattern(root *SingleDirV2, dir string, pattern string) {
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
