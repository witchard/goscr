# goscr

Use go like its a scripting language

Inspired by https://github.com/bitfield/script, this project aims to make it even easier to write go "scripts". The aim is to eventually be able to do something like `goscr -C 'script.Stdin().Match("Error").Stdout()'` to quickly run a go script.

## Installation

You need `goimports` and the `go` compiler installed and on your `$PATH`:

- Get go from https://go.dev
- Install `goimports` with `go install golang.org/x/tools/cmd/goimports@latest`

Then install this with:

- `go install github.com/witchard/goscr@latest`

## Usage

Simply pass a file with your go code in... this code is wrapped up as a function and executed. An error called `err` is already defined as a [named return value](https://go.dev/tour/basics/7) - set this to non-nil if you want to exit with an error. Because you are in a function, you can just `return` to exit your script early.

Under the hood, a `.goscr` directory in your home drive is used to compile your scripts as full go programs.

## To do

This project has only just started... we still need to:

- [] Support compiling the same program twice - currently bombs out
- [] Used cached compiled code when the same script is run
- [] Lock compilation directory so that parallel runs don't interfere with each other
- [] Support command line options for your scripts
- [] Support `-C` for passing code on the command line
- [] Allow user to hint at what imports are needed
- [] Provide better mechanism for presenting compilation errors back to the user (map line numbers)
- [] Clean up old compilation directories
- [] Test on mac and windows
- [] Support mode where code is run against each line of stdin
- [] Support extra user code from different files
- [] ?Add support for testing scripts?
