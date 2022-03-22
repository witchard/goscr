package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

var logger *log.Logger

func main() {
	args := flag.NewFlagSet("goscr", flag.ContinueOnError)
	var debug bool
	var code string
	args.BoolVar(&debug, "d", false, "Enable debug logging")
	args.StringVar(&code, "c", "", "Pass code on the command line instead of script file")
	if err := args.Parse(os.Args[1:]); err != nil {
		fmt.Println("Pass script file (if not using -C) and script args after the above flags")
		os.Exit(1)
	}

	if debug {
		logger = log.New(os.Stdout, "goscr", log.LstdFlags)
	} else {
		logger = log.New(io.Discard, "goscr", log.LstdFlags)
	}

	var runArgs []string
	if code != "" {
		runArgs = append([]string{"goscr"}, args.Args()...)
	} else {
		if len(args.Args()) == 0 {
			args.Usage()
			fmt.Println("Pass script file (if not using -C) and script args after the above flags")
			os.Exit(1)
		}

		codeRaw, err := os.ReadFile(args.Arg(0))
		if err != nil {
			log.Fatalln(err)
		}
		code = string(codeRaw)

		runArgs = args.Args()
	}

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
	err = syscall.Exec(filepath.Join(workdir, "goscr"), runArgs, os.Environ())
	if err != nil {
		log.Fatalln("Failed to execute compiled script in", workdir, err)
	}
}
