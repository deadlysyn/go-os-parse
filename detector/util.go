package detector

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
)

// Try to read os-release file.
// https://www.freedesktop.org/software/systemd/man/os-release.html
func readReleaseFile() ([]byte, error) {
	var data []byte

	// Supported locations, in order of precedence.
	files := []string{
		"/etc/os-release",
		"/usr/lib/os-release",
	}

	var file string
	for _, f := range files {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			// Must stop at first match...
			file = f
			break
		}
	}

	if file == "" {
		return data, errors.New("Unable to read os-release")
	}

	f, err := os.Open(file)
	if err != nil {
		return data, err
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return data, err
	}

	data = make([]byte, s.Size())
	_, err = f.Read(data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func parseField(raw [][]byte, field string) ([]byte, error) {
	var parsed []byte

	for _, v := range raw {
		if matched, _ := regexp.Match(fmt.Sprintf("^%s=.*$", field), v); matched {
			parsed = bytes.Split(v, []byte("="))[1]
			parsed = bytes.Trim(parsed, "\" ")
			return bytes.ToLower(parsed), nil
		}
	}
	return parsed, fmt.Errorf("Failed to parse %s field", field)
}
