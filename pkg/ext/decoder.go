package ext

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

// SingleDir struct
type SingleDir struct {
	Name              string       `yaml:"name"`
	Path              string       `yaml:"path"`
	Files             *[]string    `yaml:"files"`
	Dirs              *[]SingleDir `yaml:"dirs"`
	AllowedPatterns   *[]string    `yaml:"allowedPatterns"`
	ForbiddenPatterns *[]string    `yaml:"forbiddenPatterns"`
	MinCount          int          `default:"0" yaml:"minCount"`
	MaxCount          int          `default:"1000" yaml:"maxCount"`
	MaxDepth          int          `default:"1000" yaml:"maxDepth"`
	AllowChildren     bool         `default:"true" yaml:"allowChildren"`
}

// UnmarshalYAML implements yaml.Unmarshaler interface
// Meant for initializing default values
func (s *SingleDir) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := defaults.Set(s)
	if err != nil {
		fmt.Println("Cannot set defaults!")
		os.Exit(1)
	}

	type plain SingleDir
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

type config struct {
	Name      string     `yaml:"name"`
	Structure *SingleDir `yaml:"structure"`
}

func ReadConfig(path string) (config, []SingleDir) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		// add emoji
		PrintOut("☠️ Can not read file ==>", path)
		os.Exit(1)
	}
	PrintOut("✅ Loaded template from  ==>", path)
	// Map to store the parsed YAML data
	var curConfig config

	// Unmarshal the YAML string into the data map
	err = yaml.Unmarshal([]byte(yamlFile), &curConfig)
	if err != nil {
		PrintOut("☠️ Can not decode file ==>", path)
		os.Exit(1)
	}

	if curConfig.Structure == nil {
		PrintOut("☠️ Can not find correct structure in ==>", path)
		os.Exit(1)
	} else {
		PrintOut("✅ Read structure for ==>", curConfig.Name)
	}

	dirList := []SingleDir{}
	traverseStructure(curConfig.Structure, ".", &dirList)
	return curConfig, dirList
}

func traverseStructure(dir *SingleDir, path string, dirsList *[]SingleDir) {
	// Change path
	if dir == nil {
		PrintOut("☠️  Can't traverse nil dir ==>", path)
		os.Exit(1)
	}
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
