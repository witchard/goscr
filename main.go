package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var dbg *log.Logger

func main() {
	args := flag.NewFlagSet("goscr", flag.ContinueOnError)
	var debug bool
	var keep bool
	var force bool
	var code string
	var imports []string
	args.BoolVar(&debug, "d", false, "Enable debug logging")
	args.BoolVar(&keep, "k", false, "Keep temporary files even on compilation error")
	args.BoolVar(&force, "f", false, "Force rebuild even if code is already compiled")
	args.StringVar(&code, "c", "", "Pass code on the command line instead of script file")
	args.Func("i", "Import hint (can specify multiple times)", func(i string) error {
		imports = append(imports, i)
		return nil
	})
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
			fmt.Println("Pass script file (if not using -c) and script args after the above flags")
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

	hash, err := HashAndCreateIfNeeded(code, keep, force, imports)
	if err != nil {
		log.Fatalln(err)
	}

	err = RunProgram(hash, runArgs)
	if err != nil {
		log.Fatalln(err)
	}
}

func HashAndCreateIfNeeded(code string, keep, force bool, imports []string) (string, error) {
	// Check if already compiled
	hash := Hash(code)
	dbg.Println("Code hash is", hash)

	// Grab read lock for check
	rd, err := LockRead(hash)
	if err != nil {
		return "", err
	}
	defer rd.Unlock()

	basedir, err := Workdir()
	if err != nil {
		return "", err
	}
	workdir := filepath.Join(basedir, hash)
	exists, err := DirExists(workdir)
	if err != nil {
		return "", err
	}

	if force || !exists {
		// Upgrade to write lock
		rd.Unlock()
		wr, err := LockWrite(hash)
		if err != nil {
			return "", err
		}
		defer wr.Unlock()

		if force && exists {
			dbg.Println("Removing existing dir", workdir)
			os.RemoveAll(workdir)
		}

		dbg.Println("Creating code in", workdir, "with imports", imports)
		err = Create(code, imports, workdir)
		if err != nil {
			os.RemoveAll(workdir) // Cleanup as something failed
			return "", fmt.Errorf("failed to create code in %s: %s", workdir, err)
		}

		dbg.Println("Compiling code")
		err = Compile(workdir)
		if err != nil {
			if !keep {
				os.RemoveAll(workdir) // Cleanup as something failed
			}
			return "", fmt.Errorf("failed to compile code in %s: %s", workdir, err)
		}
	}
	return hash, nil
}

func RunProgram(hash string, args []string) error {
	lck, err := LockRead(hash)
	if err != nil {
		return err
	}
	defer lck.Unlock()

	basedir, err := Workdir()
	if err != nil {
		return err
	}
	workdir := filepath.Join(basedir, hash)

	binary := filepath.Join(workdir, "goscr")
	dbg.Println("Executing", binary, "with args", args)
	cmd := exec.Command(binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute compiled script in %s: %s", workdir, err)
	}
	return nil
}
