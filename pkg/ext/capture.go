package ext

// CaptureOptions controls a capture run.
type CaptureOptions struct {
	Name       string
	OutputPath string
	Depth      int
	Root       string // defaults to "." when empty
}

// CaptureReport summarizes a Capture run.
type CaptureReport struct {
	Dirs     int
	Files    int
	Patterns int
}

// captureScan holds the raw filesystem scan results before they are
// rendered to a config file. Kept unexported because the only legit
// consumer is Capture itself.
type captureScan struct {
	report   CaptureReport
	subdirs  []string
	files    []string
	patterns map[string]string
}

// Capture walks the filesystem at opts.Root, derives a structure
// template, and writes it to opts.OutputPath. Equivalent to calling
// CaptureScan followed by ExportConfig.
func Capture(opts CaptureOptions) (CaptureReport, error) {
	s := scan(opts)
	ExportConfig(opts.Name, opts.OutputPath, s.subdirs, s.files, s.patterns)
	return s.report, nil
}

// CaptureWithReporter is like Capture but invokes onScan after the
// filesystem walk and before writing the YAML. It exists so the CLI
// can print the dir/file/pattern counts in the same order users have
// always seen them.
func CaptureWithReporter(opts CaptureOptions, onScan func(CaptureReport)) (CaptureReport, error) {
	s := scan(opts)
	if onScan != nil {
		onScan(s.report)
	}
	ExportConfig(opts.Name, opts.OutputPath, s.subdirs, s.files, s.patterns)
	return s.report, nil
}

func scan(opts CaptureOptions) captureScan {
	root := opts.Root
	if root == "" {
		root = "."
	}
	subdirs := CaptureSubDirs(root, opts.Depth)
	files := CaptureAllFiles(root, opts.Depth)
	patterns := CaptureRequiredPatterns(root, opts.Depth)
	return captureScan{
		report: CaptureReport{
			Dirs:     len(subdirs),
			Files:    len(files),
			Patterns: len(patterns),
		},
		subdirs:  subdirs,
		files:    files,
		patterns: patterns,
	}
}
