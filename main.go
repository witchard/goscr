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

	workdir, err := Create(os.Args[1], []string{})
	if err != nil {
		log.Fatalln("Failed to create script compilation directory", err)
	}

	err = Compile(workdir)
	if err != nil {
		log.Fatalln("Failed to compile script in", workdir, err)
	}

	err = syscall.Exec(filepath.Join(workdir, "goscr"), []string{}, os.Environ())
	if err != nil {
		log.Fatalln("Failed to execute compiled script in", workdir, err)
	}
}
