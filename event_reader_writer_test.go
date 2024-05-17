package main

import (
	"bufio"
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
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(c.InputPath)
			if err != nil {
				t.Fatal(err)
			}
			outF, err := os.Open(c.OutputPath)
			defer func() {
				cerr := f.Close()
				if cerr != nil {
					err = cerr
				}
			}()
			defer func() {
				cerr := outF.Close()
				if cerr != nil {
					err = cerr
				}
			}()

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
	t.Strings = append(t.Strings, strings.SplitAfter(string(p), "\n")...)
	return len(p), nil
}
