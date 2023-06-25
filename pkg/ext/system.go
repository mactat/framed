package ext

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CreateDir(path string) {
	// Create directory
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		fmt.Printf("%-35s %-35s\n", "📂 Creating directory ==> ", path)
		if err != nil {
			log.Println(err)
		}
	}

}

func CreateAllDirs(dirList []SingleDir) {
	for _, dir := range dirList {
		CreateDir(dir.Path)
	}
}

func CreateFile(path string, name string) {
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

func CreateAllFiles(dirList []SingleDir) {
	for _, dir := range dirList {
		if dir.Files == nil {
			continue
		}
		for _, file := range *dir.Files {
			CreateFile(dir.Path, file)
		}
	}
}

// Check if directory exists on given path and is type dir
func DirExists(path string) bool {
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
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Count files in given directory
func CountFiles(path string) int {
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

func HasDirs(path string) bool {
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
func CheckDepth(path string) int {
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

func MatchPatternInDir(path string, pattern string) []string {
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
func CaptureSubDirs(path string, depth int) []string {
	var dirs []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		} else if info.IsDir() && depth > 0 && strings.Count(path, string(os.PathSeparator)) >= depth {
			return filepath.SkipDir
		} else if info.IsDir() && info.Name() != "." {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs
}

// Capture all files in given directory
func CaptureAllFiles(path string, depth int) []string {
	var files []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		} else if !info.IsDir() && depth > 0 && strings.Count(path, string(os.PathSeparator)) >= depth {
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
func CaptureRequiredPatterns(path string, depth int) map[string]string {
	var rules = make(map[string]string)
	var dirs []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		} else if info.IsDir() && depth > 0 && strings.Count(path, string(os.PathSeparator)) >= depth {
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

func VerifyFiles(dir SingleDir, allGood *bool) {
	for _, file := range *dir.Files {
		if !FileExists(dir.Path + "/" + file) {
			PrintOut("❌ File not found      ==>", dir.Path+"/"+file)
			*allGood = false
		}
	}
}
func VerifyForbiddenPatterns(dir SingleDir, allGood *bool) {
	for _, pattern := range *dir.ForbiddenPatterns {
		matched := MatchPatternInDir(dir.Path, pattern)
		for _, match := range matched {
			PrintOut("❌ Forbidden pattern ("+pattern+") matched under ==>", dir.Path+"/"+match)
			*allGood = false
		}
	}
}

func VerifyAllowedPatterns(dir SingleDir, allGood *bool) {
	matchedCount := 0
	for _, pattern := range *dir.AllowedPatterns {
		matched := MatchPatternInDir(dir.Path, pattern)
		matchedCount += len(matched)
	}
	if matchedCount != CountFiles(dir.Path) && len(*dir.AllowedPatterns) > 0 {
		patternsString := strings.Join(*dir.AllowedPatterns, " ")
		PrintOut("❌ Not all files match required pattern ("+patternsString+") under ==>", dir.Path)
		*allGood = false
	}
}
