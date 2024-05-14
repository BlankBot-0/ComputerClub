package main

import (
	"bufio"
	"os"
	"testing"
)

func TestEventReaderWriter_BasicScenario(t *testing.T) {
	f, err := os.Open("./test_file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	in := bufio.NewReader(f)
	out := TestWriter{
		Strings: make([]string, 0),
	}
	err = Handle(in, out)
	if err != nil {
		t.Fatal(err)
	}

}

type TestWriter struct {
	Strings []string
}

func (t TestWriter) Write(p []byte) (int, error) {
	t.Strings = append(t.Strings, string(p))
	return len(t.Strings), nil
}
