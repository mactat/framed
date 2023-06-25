package ext

import (
	"io"
	"net/http"
	"os"
)

func ExampleToUrl(example string) string {
	return "https://raw.githubusercontent.com/mactat/framed/master/examples/" + example + ".yaml"
}

func ImportFromUrl(path string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
