package ext

// CreateOptions controls what Create produces.
type CreateOptions struct {
	TemplatePath string
	CreateFiles  bool
}

// CreateReport summarizes a Create run.
type CreateReport struct {
	DirsCreated  int
	FilesCreated int
}

// Create reads a template and materializes the dirs (and optionally
// the files) it defines, relative to the current working directory.
func Create(opts CreateOptions) (CreateReport, error) {
	_, dirList, err := ReadConfigE(opts.TemplatePath)
	if err != nil {
		return CreateReport{}, err
	}

	CreateAllDirs(dirList)

	report := CreateReport{DirsCreated: len(dirList)}
	if opts.CreateFiles {
		CreateAllFiles(dirList)
		for _, dir := range dirList {
			if dir.Files != nil {
				report.FilesCreated += len(*dir.Files)
			}
		}
	}
	return report, nil
}
