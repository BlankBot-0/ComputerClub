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
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()

	in := bufio.NewReader(f)
	err = Handle(in, os.Stdout)
	if err != nil {
		panic(err)
	}
}
