package ext

import "fmt"

// ImportOptions controls an import run. Exactly one of URL or Example
// should be set — that mirrors the Cobra MarkFlagsMutuallyExclusive
// in the cmd layer.
type ImportOptions struct {
	URL        string
	Example    string
	OutputPath string
}

// ImportResult describes a successful import.
type ImportResult struct {
	Name   string
	Source string
}

// Import downloads a template from a URL or from the built-in
// examples bucket, writes it to OutputPath, and parses it to confirm
// it's a valid template.
func Import(opts ImportOptions) (ImportResult, error) {
	var source string
	switch {
	case opts.URL != "" && opts.Example != "":
		return ImportResult{}, fmt.Errorf("provide either URL or Example, not both")
	case opts.URL != "":
		source = opts.URL
	case opts.Example != "":
		source = ExampleToUrl(opts.Example)
	default:
		return ImportResult{}, fmt.Errorf("either URL or Example must be set")
	}

	if err := ImportFromUrl(opts.OutputPath, source); err != nil {
		return ImportResult{Source: source}, fmt.Errorf("download from %s: %w", source, err)
	}

	cfg, _, err := ReadConfigE(opts.OutputPath)
	if err != nil {
		return ImportResult{Source: source}, err
	}
	return ImportResult{Name: cfg.Name, Source: source}, nil
}
