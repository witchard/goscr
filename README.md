# goscr

Use go like its a scripting language.

Inspired by https://github.com/bitfield/script, this project aims to make it even easier to write go "scripts". You can for example do things like `goscr -c 'script.Stdin().Match("Error").Stdout()'` to quickly run a go script. You can also pass a script in as a file as `goscr <somescript>`, or use `goscr` in a shebang.

## Installation

You the `go` compiler installed, get it from https://go.dev.

Then install this with:

- `go install github.com/witchard/goscr@latest`

`goscr` should work on windows, linux and mac.

## Usage

Simply pass a file with your go code in... this code is wrapped up as a function and executed. An error called `err` is already defined as a [named return value](https://go.dev/tour/basics/7) - set this to non-nil if you want to exit with an error. Because you are in a function, you can just `return` to exit your script early.

Under the hood, a `.goscr` directory in your home drive is used to compile your scripts as full go programs.

A function `P` is available that behaves like `fmt.Println` for printing stuff out. A similar function `E` is also available, but prints to stderr.

A function `L` is available for processing each line of stdin, it calls a callback function for each line (simple data types are converted automatically), for example:

```bash
echo -e "1\n2\n3" | goscr -c "s := 0; err = L(func(i int){s += i}); P(s)"
```

`goscr` also works with a shebang (`#!`), see [example.goscr](example.goscr).

## To do

This project has only just started... we still need to:

- [X] Support compiling the same program twice - currently bombs out
- [X] Used cached compiled code when the same script is run
- [X] Lock compilation directory so that parallel runs don't interfere with each other
- [X] Support command line options for your scripts
- [X] Support `-c` for passing code on the command line
- [X] Allow user to hint at what imports are needed
- [X] Support working with `#!`
- [ ] Provide better mechanism for presenting compilation errors back to the user (map line numbers)
- [ ] Clean up old compilation directories
- [X] Test on mac and windows
- [X] Support mode where code is run against each line of stdin
- [ ] Improve module documentation (and of lines)
- [ ] Add more unit tests :-)

## Alternatives

* For the ability to run go programs with a shebang: https://github.com/erning/gorun
* Extended go syntax with scripting feel: https://github.com/goplus/gop