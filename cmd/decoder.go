/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/ioutil"
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

// SingleDir struct
type SingleDir struct {
	Name              string       `yaml:"name"`
	Path              string       `yaml:"path"`
	Files             []string     `yaml:"files"`
	Dirs              *[]SingleDir `yaml:"dirs"`
	AllowedPatterns   []string     `yaml:"allowedPatterns"`
	ForbiddenPatterns []string     `yaml:"forbiddenPatterns"`
	MinCount          int          `default:"0" yaml:"minCount"`
	MaxCount          int          `default:"1000" yaml:"maxCount"`
	MaxDepth          int          `default:"1000" yaml:"maxDepth"`
	AllowChildren     bool         `default:"true" yaml:"allowChildren"`
}

// UnmarshalYAML implements yaml.Unmarshaler interface
// Meant for initializing default values
func (s *SingleDir) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(s)

	type plain SingleDir
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

type Conf struct {
	Name      string     `yaml:"name"`
	Structure *SingleDir `yaml:"structure"`
}

func readConfig(path string) (Conf, []SingleDir) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		// add emoji
		print("☠️  Can't read file ==>", path)
		os.Exit(1)
	}
	print("✅ Loaded template from  ==>", path)
	// Map to store the parsed YAML data
	var config Conf

	// Unmarshal the YAML string into the data map
	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		print("☠️  Can't decode file ==>", path)
		os.Exit(1)
	}
	print("✅ Read structure for ==>", config.Name)

	dirList := []SingleDir{}
	traverseStructure(config.Structure, ".", &dirList)
	return config, dirList
}

func traverseStructure(dir *SingleDir, path string, dirsList *[]SingleDir) {
	// Chanage path
	dir.Path = path

	// add current dir to dirsList
	*dirsList = append(*dirsList, *dir)

	if dir.Dirs == nil {
		return
	}
	// traverse children
	for _, child := range *dir.Dirs {
		traverseStructure(&child, path+"/"+child.Name, dirsList)
	}
}
