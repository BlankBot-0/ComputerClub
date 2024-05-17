package src

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestAllScenarios(t *testing.T) {
	testsPath := "./../test_files/"
	cases := []struct {
		InputPath          string
		ExpectedOutputPath string
		Name               string
	}{
		{
			InputPath:          "test_file.txt",
			ExpectedOutputPath: "output/test_output_file.txt",
			Name:               "BasicTest",
		},
		{
			InputPath:          "test_2.txt",
			ExpectedOutputPath: "output/test_2.txt",
			Name:               "Test2",
		},
		{
			InputPath:          "test_3.txt",
			ExpectedOutputPath: "output/test_3.txt",
			Name:               "Test3",
		},
		{
			InputPath:          "autoleave_test.txt",
			ExpectedOutputPath: "outputs/autoleave_test_output.txt",
			Name:               "AutoLeaveTest",
		},
		{
			InputPath:          "wrong_format_time_format_test.txt",
			ExpectedOutputPath: "outputs/wrong_format_time_format_test.txt",
			Name:               "TimeFormatZerosTest",
		},
		{
			InputPath:          "wrong_format_deskNum_test.txt",
			ExpectedOutputPath: "outputs/wrong_format_deskNum_test.txt",
			Name:               "DeskNumFormatTest",
		},
		{
			InputPath:          "wrong_format_time_travel_test.txt",
			ExpectedOutputPath: "outputs/wrong_format_time_travel_test.txt",
			Name:               "FormatTimeTravelTest",
		},
		{
			InputPath:          "wrong_format_bad_name_test.txt",
			ExpectedOutputPath: "outputs/wrong_format_bad_name_test.txt",
			Name:               "BadNameFormatTest",
		},
		{
			InputPath:          "wrong_format_missing_fields.txt",
			ExpectedOutputPath: "outputs/wrong_format_missing_fields.txt",
			Name:               "MissingFieldsFormatTest",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			f, err := os.Open(testsPath + c.InputPath)
			if err != nil {
				t.Fatal(err)
			}
			outF, err := os.Open(testsPath + c.ExpectedOutputPath)
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
	t.Strings = append(t.Strings, strings.SplitAfter(string(p), "\n")...)
	return len(p), nil
}
