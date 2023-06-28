package ext

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Consider implementing a custom Unmarshaler for SingleDirOut
type SingleDirOut struct {
	Name              string          `yaml:"name"`
	Path              string          `yaml:"-"`
	Files             *[]string       `yaml:"files,omitempty"`
	Dirs              *[]SingleDirOut `yaml:"dirs,omitempty"`
	AllowedPatterns   *[]string       `yaml:"allowedPatterns,omitempty"`
	ForbiddenPatterns *[]string       `yaml:"forbiddenPatterns,omitempty"`
	MinCount          int             `yaml:"minCount,omitempty"`
	MaxCount          int             `yaml:"maxCount,omitempty"`
	MaxDepth          int             `yaml:"maxDepth,omitempty"`
	AllowChildren     bool            `yaml:"allowChildren,omitempty"`
}

type configOut struct {
	Name      string        `yaml:"name"`
	Structure *SingleDirOut `yaml:"structure"`
}

func ExportConfig(name string, path string, subdirs []string, files []string, patterns map[string]string) {
	// create config, files and dirs are empty
	config := configOut{
		Name: name,
		Structure: &SingleDirOut{
			Name:  "root",
			Files: &[]string{},
			Dirs:  &[]SingleDirOut{},
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
	err = os.WriteFile(path, yamlFile, 0644)
	if err != nil {
		fmt.Printf("Error while writing file. %v", err)
	}
}

func insertSubdirs(dirs *[]SingleDirOut, subdirs []string) {
	for subdir := range subdirs {
		insertSingleDir(dirs, subdirs[subdir])
	}
}

func insertSingleDir(dirs *[]SingleDirOut, dir string) {
	subdirPath := strings.Split(dir, "/")
	curDirs := dirs
	for i := range subdirPath {
		if !containsDir(*curDirs, subdirPath[i]) {
			*curDirs = append(*curDirs, SingleDirOut{
				Name:            subdirPath[i],
				Files:           &[]string{},
				Dirs:            &[]SingleDirOut{},
				AllowedPatterns: &[]string{},
			})
		}
		// go deeper
		curDirs = getDir(*curDirs, subdirPath[i]).Dirs
	}
}

func containsDir(dirs []SingleDirOut, name string) bool {
	for _, dir := range dirs {
		if dir.Name == name {
			return true
		}
	}
	return false
}

func getDir(dirs []SingleDirOut, name string) *SingleDirOut {
	for _, dir := range dirs {
		if dir.Name == name {
			return &dir
		}
	}
	return nil
}

func insertFiles(root *SingleDirOut, files []string) {
	for file := range files {
		insertSingleFile(root, files[file])
	}
}

func insertSingleFile(root *SingleDirOut, file string) {
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

func insertPatterns(root *SingleDirOut, patterns map[string]string) {
	for dir, pattern := range patterns {
		insertSinglePattern(root, dir, pattern)
	}
}

func insertSinglePattern(root *SingleDirOut, dir string, pattern string) {
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
