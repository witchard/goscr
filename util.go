package main

import (
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"os"
	"path/filepath"
)

func Hash(code string) string {
	sha := sha1.Sum([]byte(code))
	return base32.StdEncoding.EncodeToString(sha[:])
}

func Workdir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".goscr"), nil
}

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
