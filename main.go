package main

import (
	"bufio"
	"os"
	"yadroTest/src"
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
	err = src.Handle(in, os.Stdout)
	if err != nil {
		panic(err)
	}
}
