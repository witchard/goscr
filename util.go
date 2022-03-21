package main

import (
	"os"
	"path/filepath"
)

func Workdir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".goscr"), nil
}
