package ext

import "fmt"

// Failure is a single rule violation found by Verify. The two-field
// shape mirrors PrintOut so cmd can render it without reformatting.
type Failure struct {
	Prompt string
	Detail string
}

// VerifyReport describes the outcome of a Verify run.
type VerifyReport struct {
	Failures []Failure
	OK       bool
}

// Verify loads a template and checks every dir against the rules
// (existence, file presence, min/max counts, depth, allowed/forbidden
// patterns, child policy). It never calls os.Exit — orchestration owns
// process lifecycle. A non-nil error means the template itself could
// not be loaded; a returned report with OK=false means rules failed.
func Verify(templatePath string) (VerifyReport, error) {
	_, dirList, err := ReadConfigE(templatePath)
	if err != nil {
		return VerifyReport{}, err
	}

	report := VerifyReport{OK: true}
	for _, dir := range dirList {
		verifyDir(dir, &report)
	}
	return report, nil
}

func verifyDir(dir SingleDir, report *VerifyReport) {
	if !DirExists(dir.Path) {
		report.fail("❌ Directory not found ==>", dir.Path)
		// Skip the rest — the disk-walking helpers (CountFiles,
		// HasDirs, CheckDepth, MatchPatternInDir) panic on a
		// missing dir, and we already know the dir is wrong.
		return
	}

	if dir.Files != nil {
		for _, file := range *dir.Files {
			if !FileExists(dir.Path + "/" + file) {
				report.fail("❌ File not found      ==>", dir.Path+"/"+file)
			}
		}
	}

	numFiles := CountFiles(dir.Path)
	if numFiles < dir.MinCount {
		report.fail("❌ Min count ("+fmt.Sprint(dir.MinCount)+") not met ==>", dir.Path)
	}
	if numFiles > dir.MaxCount {
		report.fail("❌ Max count ("+fmt.Sprint(dir.MaxCount)+") exceeded ==>", dir.Path)
	}

	if !dir.AllowChildren && HasDirs(dir.Path) {
		report.fail("❌ Children not allowed ==>", dir.Path)
	}

	if CheckDepth(dir.Path) > dir.MaxDepth {
		report.fail("❌ Max depth exceeded ("+fmt.Sprint(dir.MaxDepth)+") ==>", dir.Path)
	}

	if dir.ForbiddenPatterns != nil {
		for _, pattern := range *dir.ForbiddenPatterns {
			for _, match := range MatchPatternInDir(dir.Path, pattern) {
				report.fail("❌ Forbidden pattern ("+pattern+") matched under ==>", dir.Path+"/"+match)
			}
		}
	}

	if dir.AllowedPatterns != nil && len(*dir.AllowedPatterns) > 0 {
		matchedCount := 0
		for _, pattern := range *dir.AllowedPatterns {
			matchedCount += len(MatchPatternInDir(dir.Path, pattern))
		}
		if matchedCount != CountFiles(dir.Path) {
			patternsString := ""
			for i, p := range *dir.AllowedPatterns {
				if i > 0 {
					patternsString += " "
				}
				patternsString += p
			}
			report.fail("❌ Not all files match required pattern ("+patternsString+") under ==>", dir.Path)
		}
	}
}

func (r *VerifyReport) fail(prompt, detail string) {
	r.OK = false
	r.Failures = append(r.Failures, Failure{Prompt: prompt, Detail: detail})
}
