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
