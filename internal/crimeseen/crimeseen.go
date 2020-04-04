// Package crimeseen contains common paths and utility functions that come in
// handy across the codebase.
package crimeseen

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Mkdirp creates the specified directory path if it doesn't already exist.
func Mkdirp(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil && !IsFileExistsError(err) {
		return err
	}

	return nil
}

// IsFileExistsError returns true if the specified error is due to a file already
// existing on the filesystem.
func IsFileExistsError(err error) bool {
	return strings.Contains(err.Error(), "file exists")
}

// FileExists returns true if the specified file path exists.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// PaddedNumberString takes an integer value and returns a string value with a
// "0" prefix if the value is less than 10.
func PaddedNumberString(value int) string {
	numString := strconv.Itoa(value)
	if value < 10 {
		return "0" + numString
	}
	return numString
}

// WriteJSONFile writes the specified contents as JSON to the specified path.
func WriteJSONFile(path string, contents interface{}) error {
	b, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, string(b))
	if err != nil {
		return err
	}

	return file.Sync()
}

// RunCommand is a wrapper around exec.Command that redirects output to the
// terminal.
func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
