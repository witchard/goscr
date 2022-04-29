package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	format "golang.org/x/tools/imports"
)

const header string = `package main

import (
	"fmt"
	"os"
	"github.com/witchard/goscr/lines"
`

const body string = `)

func P(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func E(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}

func L(cb interface{}) error {
	return lines.EachStdin(cb)
}

func main() {
	if err := __run_goscr__(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func __run_goscr__() (err error) {
	// ---- START USER CODE ----
`

const footer string = `
	// ---- END USER CODE ----
	return
}
`

func Create(code string, imports []string, workdir string) error {
	// Create dir
	err := os.MkdirAll(workdir, 0o700)
	if err != nil {
		return err
	}

	// Setup go.mod
	if err := Run(workdir, "go", "mod", "init", "goscr"); err != nil {
		return err
	}

	// Build file contents
	buf := &bytes.Buffer{}

	_, err = fmt.Fprint(buf, header)
	if err != nil {
		return err
	}

	for _, dep := range imports {
		_, err = fmt.Fprintf(buf, "\t\"%s\"\n", dep)
		if err != nil {
			return err
		}

	}

	_, err = fmt.Fprint(buf, body)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(buf, code)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(buf, footer)
	if err != nil {
		return err
	}

	// Sort imports / format file
	mainFile := filepath.Join(workdir, "main.go")
	formatted, err := format.Process(mainFile, buf.Bytes(), nil)
	if err != nil {
		return err
	}

	// Write file
	file, err := os.Create(mainFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(formatted)
	return err
}
