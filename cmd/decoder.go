/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Required struct {
	Patterns []string `mapstructure:"patterns"`
	Files    []string `mapstructure:"files"`
}

type Forbidden struct {
	Patterns []string `mapstructure:"patterns"`
	Files    []string `mapstructure:"files"`
}

// SingleDir struct
type SingleDir struct {
	Path          string    `mapstructure:"path"`
	Required      Required  `mapstructure:"required"`
	Forbidden     Forbidden `mapstructure:"forbidden"`
	MinCount      int       `mapstructure:"minCount"`
	MaxCount      int       `mapstructure:"maxCount"`
	MaxDepth      int       `mapstructure:"maxDepth"`
	AllowChildren bool      `mapstructure:"allowChildren"`
}

func readConfig(path string) []SingleDir {
	config := readYaml(path)
	name := config.(map[string]interface{})["name"]
	nameString := fmt.Sprintf("%v", name)
	print("✅ Loading template from  ==>", path)
	print("✅ Reading structure for ==>", nameString+"\n")
	dirsList := []SingleDir{}

	// traverse the structure to flatten it
	traverseStructure(config.(map[string]interface{})["structure"].(map[string]interface{})["root"], ".", &dirsList)
	return dirsList
}

func readYaml(path string) interface{} {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		// add emoji
		print("☠️  Can't read file ==>", path)
		os.Exit(1)
	}
	// Map to store the parsed YAML data
	var data map[string]interface{}

	// Unmarshal the YAML string into the data map
	err = yaml.Unmarshal([]byte(yamlFile), &data)
	if err != nil {
		print("☠️  Can't decode file ==>", path)
		os.Exit(1)
	}
	return data
}

func newSingleDir() SingleDir {
	singleDir := SingleDir{}
	// set default values
	singleDir.MaxDepth = 1000
	singleDir.MaxCount = 1000
	singleDir.MinCount = 0
	singleDir.AllowChildren = true
	return singleDir
}

func decodeSingleDir(data interface{}) SingleDir {
	decoded := newSingleDir()

	// decode
	err := mapstructure.Decode(data, &decoded)
	if err != nil {
		fmt.Println("☠️  Wrong structure in template file")
		os.Exit(1)
	}
	return decoded
}

func traverseStructure(data interface{}, path string, dirsList *[]SingleDir) {
	// decode current dir
	single := decodeSingleDir(data)
	single.Path = path
	*dirsList = append(*dirsList, single)

	// dir def without any properties by checking the type
	if _, ok := data.(map[string]interface{}); !ok {
		return
	}
	// decode children
	children, ok := data.(map[string]interface{})["children"]
	if !ok {
		return
	}
	for _, value := range children.([]interface{}) {
		for name, body := range value.(map[string]interface{}) {
			traverseStructure(body, path+"/"+name, dirsList)
		}
	}
}
