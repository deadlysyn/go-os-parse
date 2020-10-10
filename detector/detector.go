package detector

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

// PackageManagerCmd attempts to return the appropriate pacakge
// manager command based on the os-release distribution ID.
// Reliance on external commands, shells or builtins is avoided.
func PackageManagerCmd() (string, error) {
	var cmd string

	data, err := readReleaseFile()
	if err != nil {
		return cmd, err
	}
	raw := bytes.Split(data, []byte("\n"))

	id, err := parseField(raw, "ID")
	if err != nil {
		return cmd, err
	}

	switch string(id) {
	case "alpine":
		cmd = "apk"
	case "arch", "manjaro":
		cmd = "pacman"
	case "debian", "ubuntu", "mint", "linspire":
		cmd = "apt"
	case "suse":
		cmd = "yast"
	case "gentoo":
		cmd = "emerge"
	case "centos", "fedora", "rhel", "redhat":
		version, err := parseField(raw, "VERSION_ID")
		if err != nil {
			return cmd, err
		}
		if v, _ := strconv.Atoi(string(version)); v <= 7 {
			cmd = "yum"
		} else {
			cmd = "dnf"
		}
	default:
		return cmd, errors.New("Failed to detect package manager")
	}

	return cmd, nil
}

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
		if bytes.HasPrefix(v, []byte(fmt.Sprintf("%s=", field))) {
			parsed = bytes.Split(v, []byte("="))[1]
			return bytes.ToLower(parsed), nil
		}
	}
	return parsed, fmt.Errorf("Failed to parse %s field", field)
}
