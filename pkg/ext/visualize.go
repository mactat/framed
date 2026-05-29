package ext

// Visualize loads a template and prints its dir/file tree.
func Visualize(templatePath string) error {
	_, dirList, err := ReadConfigE(templatePath)
	if err != nil {
		return err
	}
	VisualizeTemplate(dirList)
	return nil
}
