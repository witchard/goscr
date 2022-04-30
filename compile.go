package main

// Compile handles the compilation of a script that has been written into workdir
func Compile(workdir string) error {
	if err := Run(workdir, "go", "mod", "tidy"); err != nil {
		return err
	}
	if err := Run(workdir, "go", "build"); err != nil {
		return err
	}
	return nil
}
