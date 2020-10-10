package detector

import (
	"bytes"
	"errors"
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
