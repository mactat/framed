package ext

// resolveImports walks the structure tree and replaces every node that
// carries a `template:` reference with the contents of the referenced
// template.
//
// Stubbed until Task 01 lands. Acts as a no-op so the decoder can call
// it unconditionally without changing existing semantics.
func resolveImports(node *SingleDir, basePath string, seen map[string]bool) error {
	return nil
}
