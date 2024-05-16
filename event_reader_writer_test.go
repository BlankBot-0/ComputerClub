package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
)

func TestAllScenarios(t *testing.T) {
	cases := []struct {
		InputPath  string
		OutputPath string
		Name       string
	}{
		{
			InputPath:  "./test_file.txt",
			OutputPath: "./test_output_file.txt",
			Name:       "BasicTest",
		},
		{
			InputPath:  "./test_files/test_2.txt",
			OutputPath: "./test_files/output/test_2.txt",
			Name:       "Test2",
		},
		{
			InputPath:  "./test_files/test_3.txt",
			OutputPath: "./test_files/output/test_3.txt",
			Name:       "Test3",
		},
		{
			InputPath:  "./test_files/autoleave_test.txt",
			OutputPath: "./test_files/outputs/autoleave_test_output.txt",
			Name:       "AutoLeaveTest",
		},
		{
			InputPath:  "./test_files/wrong_format_time_format_test.txt",
			OutputPath: "./test_files/outputs/wrong_format_time_format_test.txt",
			Name:       "TimeFormatZerosTest",
		},
		{
			InputPath:  "./test_files/wrong_format_deskNum_test.txt",
			OutputPath: "./test_files/outputs/wrong_format_deskNum_test.txt",
			Name:       "DeskNumFormatTest",
		},
		{
			InputPath:  "./test_files/wrong_format_time_travel_test.txt",
			OutputPath: "./test_files/outputs/wrong_format_time_travel_test.txt",
			Name:       "FormatTimeTravelTest",
		},
		{
			InputPath:  "./test_files/wrong_format_bad_name_test.txt",
			OutputPath: "./test_files/outputs/wrong_format_bad_name_test.txt",
			Name:       "BadNameFormatTest",
		},
		{
			InputPath:  "./test_files/wrong_format_missing_fields.txt",
			OutputPath: "./test_files/outputs/wrong_format_missing_fields.txt",
			Name:       "MissingFieldsFormatTest",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(c.InputPath)
			if err != nil {
				t.Fatal(err)
			}
			outF, err := os.Open(c.OutputPath)
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
		})
	}
}

type TestWriter struct {
	Strings []string
}

func (t *TestWriter) Write(p []byte) (int, error) {
	for _, s := range strings.SplitAfter(string(p), "\n") {
		t.Strings = append(t.Strings, s)
	}
	return len(p), nil
}

type TestReader struct {
	Strings []string
	ReadInd int
}

func (t *TestReader) Read(p []byte) (int, error) {
	if t.ReadInd >= len(t.Strings) {
		return 0, io.EOF
	}
	p = []byte(t.Strings[t.ReadInd])
	t.ReadInd++
	return len(p), nil
}
