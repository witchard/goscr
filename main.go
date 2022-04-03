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

var dbg *log.Logger

func main() {
	args := flag.NewFlagSet("goscr", flag.ContinueOnError)
	var debug bool
	var keep bool
	var code string
	args.BoolVar(&debug, "d", false, "Enable debug logging")
	args.BoolVar(&keep, "k", false, "Keep temporary files even on compilation error")
	args.StringVar(&code, "c", "", "Pass code on the command line instead of script file")
	if err := args.Parse(os.Args[1:]); err != nil {
		fmt.Println("Pass script file (if not using -C) and script args after the above flags")
		os.Exit(1)
	}

	if debug {
		dbg = log.New(os.Stderr, "debug ", log.LstdFlags)
	} else {
		dbg = log.New(io.Discard, "debug ", log.LstdFlags)
	}

	var runArgs []string
	if code != "" {
		dbg.Println("Using code from command line")
		runArgs = append([]string{"goscr"}, args.Args()...)
	} else {
		if len(args.Args()) == 0 {
			args.Usage()
			fmt.Println("Pass script file (if not using -C) and script args after the above flags")
			os.Exit(1)
		}

		dbg.Println("Reading code from", args.Arg(0))
		codeRaw, err := os.ReadFile(args.Arg(0))
		if err != nil {
			log.Fatalln(err)
		}
		code = string(codeRaw)

		runArgs = args.Args()
	}

	// Check if already compiled
	hash := Hash(code)
	dbg.Println("Code hash is", hash)
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
		dbg.Println("Creating code in", workdir)
		err = Create(code, []string{}, workdir)
		if err != nil {
			os.RemoveAll(workdir) // Cleanup as something failed
			log.Fatalln("Failed to create script compilation directory", err)
		}

		dbg.Println("Compiling code")
		err = Compile(workdir)
		if err != nil {
			if !keep {
				os.RemoveAll(workdir) // Cleanup as something failed
			}
			log.Fatalln("Failed to compile script in", workdir, err)
		}
	}

	binary := filepath.Join(workdir, "goscr")
	dbg.Println("Executing", binary, "with args", runArgs)
	err = syscall.Exec(binary, runArgs, os.Environ())
	if err != nil {
		log.Fatalln("Failed to execute compiled script in", workdir, err)
	}
}
