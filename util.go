package main

import (
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Hash provides a hash of the script code
func Hash(code string) string {
	sha := sha1.Sum([]byte(code))
	return base32.StdEncoding.EncodeToString(sha[:])
}

// Workdir computes the goscr working folder
func Workdir() (string, error) {
	if dir, ok := os.LookupEnv("GOSCR_PATH"); ok {
		return dir, os.MkdirAll(dir, 0o700)
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, ".goscr")
	return dir, os.MkdirAll(dir, 0o700)
}

// DirExists validates that the provided directory exists
func DirExists(workdir string) (bool, error) {
	stat, err := os.Stat(workdir)
	if err != nil {
		if os.IsNotExist(err) {
			// Dir does not exist
			return false, nil
		}
		return false, err // An actual error
	}
	if !stat.IsDir() {
		return false, fmt.Errorf("%s is not a directory", workdir)
	}
	return true, nil
}

// Run executes a command for use when compiling programs created from scripts
func Run(dir string, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Dir = dir
	out, err := command.Output()
	if err != nil {
		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			fmt.Println("Command", cmd, args, "exited with status", execErr.ProcessState.ExitCode())
			fmt.Println("---------- stdout ----------")
			fmt.Println(string(out))
			fmt.Println("---------- stdout ----------")
			fmt.Println(string(execErr.Stderr))
		}
	}
	return err
}

// StripShebang removes the shebang line if it is present at the beginning of the provided string
func StripShebang(in string) string {
	split := strings.SplitN(in, "\n", 2)
	if strings.HasPrefix(split[0], "#!") {
		return split[1]
	}
	return in
}
