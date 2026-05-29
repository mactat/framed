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
	Template          string       `yaml:"template"`
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

type Config struct {
	Name      string     `yaml:"name"`
	Structure *SingleDir `yaml:"structure"`
}

// config remains as a backwards-compatible alias for callers in this
// package that still use the unexported name.
type config = Config

// ReadConfigE loads a template and returns the parsed config plus the
// flattened dir list. It returns an error instead of calling os.Exit;
// the bare ReadConfig wrapper preserves the original exit-on-failure
// behavior for callers that haven't migrated yet.
func ReadConfigE(path string) (Config, []SingleDir, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		PrintOut("☠️ Can not read file ==>", path)
		return Config{}, nil, fmt.Errorf("read template %q: %w", path, err)
	}
	PrintOut("✅ Loaded template from  ==>", path)

	var curConfig Config
	if err := yaml.Unmarshal(yamlFile, &curConfig); err != nil {
		PrintOut("☠️ Can not decode file ==>", path)
		return Config{}, nil, fmt.Errorf("decode template %q: %w", path, err)
	}

	if curConfig.Structure == nil {
		PrintOut("☠️ Can not find correct structure in ==>", path)
		return Config{}, nil, fmt.Errorf("template %q has no structure", path)
	}
	PrintOut("✅ Read structure for ==>", curConfig.Name)

	if err := resolveImports(curConfig.Structure, path, map[string]bool{}); err != nil {
		return Config{}, nil, err
	}

	dirList := []SingleDir{}
	traverseStructure(curConfig.Structure, ".", &dirList)
	return curConfig, dirList, nil
}

// ReadConfig is the legacy entry point — exits the process on error.
// Prefer ReadConfigE in new code.
func ReadConfig(path string) (Config, []SingleDir) {
	cfg, dirs, err := ReadConfigE(path)
	if err != nil {
		os.Exit(1)
	}
	return cfg, dirs
}

func traverseStructure(dir *SingleDir, path string, dirsList *[]SingleDir) {
	if dir == nil {
		PrintOut("☠️  Can't traverse nil dir ==>", path)
		os.Exit(1)
	}
	dir.Path = path

	*dirsList = append(*dirsList, *dir)

	if dir.Dirs == nil {
		return
	}
	for _, child := range *dir.Dirs {
		traverseStructure(&child, path+"/"+child.Name, dirsList)
	}
}
