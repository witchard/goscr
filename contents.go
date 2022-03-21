package main

const header string = `package main

import (
`

const body string = `)

func main() {
	if err := __run_goscr__(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func __run_goscr__() (err error) {
	
`

const footer string = `
	return
}
`
