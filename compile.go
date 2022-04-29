package main

func Compile(workdir string) error {
	if err := Run(workdir, "go", "mod", "tidy"); err != nil {
		return err
	}
	if err := Run(workdir, "go", "build"); err != nil {
		return err
	}
	return nil
}
