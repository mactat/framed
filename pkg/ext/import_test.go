package ext

import "testing"

func TestImportRejectsBothFlags(t *testing.T) {
	_, err := Import(ImportOptions{URL: "x", Example: "y", OutputPath: "/tmp/out.yaml"})
	if err == nil {
		t.Fatal("expected error when both URL and Example set")
	}
}

func TestImportRejectsNoSource(t *testing.T) {
	_, err := Import(ImportOptions{OutputPath: "/tmp/out.yaml"})
	if err == nil {
		t.Fatal("expected error when neither URL nor Example set")
	}
}

func TestExampleToUrl(t *testing.T) {
	got := ExampleToUrl("golang")
	want := "https://raw.githubusercontent.com/mactat/framed/master/examples/golang.yaml"
	if got != want {
		t.Errorf("ExampleToUrl(golang) = %q, want %q", got, want)
	}
}
