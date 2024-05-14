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
	outF, err := os.Open("./test_output_file.txt")
	defer f.Close()
	defer outF.Close()

	in := bufio.NewReader(f)
	out := &TestWriter{
		Strings: make([]string, 0),
	}
	err = Handle(in, out)
	if err != nil {
		t.Fatal(err)
	}

	assertingReader := bufio.NewReader(outF)
	i := 0
	for eventStr, err := assertingReader.ReadString('\n'); err == nil; eventStr, err = assertingReader.ReadString('\n') {
		if eventStr != out.Strings[i] {
			t.Fatalf("Event at line %d: expected %s, got %s", i, eventStr, out.Strings[i])
		}
		i++
	}
}

type TestWriter struct {
	Strings []string
}

func (t *TestWriter) Write(p []byte) (int, error) {
	t.Strings = append(t.Strings, string(p))
	return len(t.Strings), nil
}
