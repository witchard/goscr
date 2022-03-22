package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const header string = `package main

import (
`

const body string = `)

func main() {
	if err := __run_goscr__(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func __run_goscr__() (err error) {
	
`

const footer string = `
	return
}
`

func Create(code string, imports []string, workdir string) error {
	// Create dir
	err := os.MkdirAll(workdir, os.ModePerm)
	if err != nil {
		return err
	}

	// Write file
	file, err := os.Create(filepath.Join(workdir, "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, header)
	if err != nil {
		return err
	}

	for _, dep := range imports {
		_, err = fmt.Fprintf(file, "\t\"%s\"\n", dep)
		if err != nil {
			return err
		}

	}

	_, err = fmt.Fprint(file, body)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(file, string(code))
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(file, footer)
	if err != nil {
		return err
	}

	return nil
}
