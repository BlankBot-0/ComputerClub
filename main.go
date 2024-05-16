package main

import (
	"bufio"
	"os"
)

func main() {
	args := os.Args[1:]

	f, err := os.Open("./" + args[0])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	in := bufio.NewReader(f)
	err = Handle(in, os.Stdout)
	if err != nil {
		panic(err)
	}
}
