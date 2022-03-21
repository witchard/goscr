package main

import (
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"os"
	"path/filepath"
)

func Create(path string, imports []string) (string, error) {
	// Slurp in the file
	code, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Compute hash of file
	sha := sha1.Sum(code)
	hash := base32.StdEncoding.EncodeToString(sha[:])

	// Create dir
	basedir, err := Workdir()
	if err != nil {
		return "", err
	}
	workdir := filepath.Join(basedir, hash)
	err = os.MkdirAll(workdir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Write file
	file, err := os.Create(filepath.Join(workdir, "main.go"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, header)
	if err != nil {
		return "", err
	}

	for _, dep := range imports {
		_, err = fmt.Fprintf(file, "\t\"%s\"\n", dep)
		if err != nil {
			return "", err
		}

	}

	_, err = fmt.Fprint(file, body)
	if err != nil {
		return "", err
	}

	_, err = fmt.Fprint(file, string(code))
	if err != nil {
		return "", err
	}

	_, err = fmt.Fprint(file, footer)
	if err != nil {
		return "", err
	}

	return workdir, nil
}
