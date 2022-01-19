package csv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsantosalmeida/csv-parser/usecase/csv"
)

func TestNewFilePatternMap(t *testing.T) {
	var (
		givenFile = "file.csv"
		input     = []byte("FName\nL Name\nS. alary\nE-mail\nID\nPhone n.\n")
		want      = map[string]*csv.FilePattern{
			givenFile: {
				FirstNameColumn: "FName",
				LastNameColumn:  "L Name",
				SalaryColumn:    "S. alary",
				EmailColumn:     "E-mail",
				IDColumn:        "ID",
				PhoneColumn:     "Phone n.",
			},
		}
	)

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed, error: %q", err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Fatalf("write args failed, error: %q", err)
	}
	w.Close()

	stdin := os.Stdin
	// Restore stdin right after the test.
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	got := csv.NewFilePatternMap([]string{givenFile})
	assert.Equal(t, want, got)
}
