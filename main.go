package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: goscr <script>")
	}

	// Slurp in the file
	codeRaw, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	code := string(codeRaw)

	// Check if already compiled
	hash := Hash(code)
	basedir, err := Workdir()
	if err != nil {
		log.Fatalln(err)
	}
	workdir := filepath.Join(basedir, hash)
	exists, err := DirExists(workdir)
	if err != nil {
		log.Fatalln(err)
	}

	if !exists {
		err = Create(code, []string{}, workdir)
		if err != nil {
			os.RemoveAll(workdir) // Cleanup as something failed
			log.Fatalln("Failed to create script compilation directory", err)
		}

		err = Compile(workdir)
		if err != nil {
			os.RemoveAll(workdir) // Cleanup as something failed
			log.Fatalln("Failed to compile script in", workdir, err)
		}
	}
	err = syscall.Exec(filepath.Join(workdir, "goscr"), []string{}, os.Environ())
	if err != nil {
		log.Fatalln("Failed to execute compiled script in", workdir, err)
	}
}
