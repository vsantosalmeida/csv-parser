package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsConstants(t *testing.T) {
	tt := []struct {
		name     string
		givenErr err
		want     string
	}{
		{
			name:     "ErrEmptyFilePatternMapReceived",
			givenErr: ErrEmptyFilePatternMapReceived,
			want:     "could not create a parser without a map of FilePattern",
		},
		{
			name:     "ErrEmptyName",
			givenErr: ErrEmptyName,
			want:     "a name is required",
		},
		{
			name:     "ErrInvalidSalaryValue",
			givenErr: ErrInvalidSalaryValue,
			want:     "could not convert salary to a float value or salary is less or equals to 0",
		},
		{
			name:     "ErrInvalidEmailFormat",
			givenErr: ErrInvalidEmailFormat,
			want:     "e-mail must be a valid address ex: email@example.com",
		},
		{
			name:     "ErrInvalidIDValue",
			givenErr: ErrInvalidIDValue,
			want:     "an id is required",
		},
		{
			name:     "ErrInvalidFilePattern",
			givenErr: ErrInvalidFilePattern,
			want:     "the columns for ID, FirstName, Salary and Email are required to process a file",
		},
		{
			name:     "ErrInvalidFilePattern",
			givenErr: ErrEmptyFilePatternMapReceived,
			want:     "could not create a parser without a map of FilePattern",
		},
		{
			name:     "ErrEmailConstraintViolation",
			givenErr: ErrEmailConstraintViolation,
			want:     "e-mail already used by an employee",
		},
		{
			name:     "ErrIDConstraintViolation",
			givenErr: ErrIDConstraintViolation,
			want:     "this ID is already used by an employee",
		},
		{
			name:     "ErrOpeningFile",
			givenErr: ErrOpeningFile,
			want:     "could not open the given file",
		},
		{
			name:     "ErrReadingFile",
			givenErr: ErrReadingFile,
			want:     "could not read the given file",
		},
		{
			name:     "ErrUnprocessableFile",
			givenErr: ErrUnprocessableFile,
			want:     "could not find a file pattern to process",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.givenErr.Error())
		})
	}
}

func TestNewError(t *testing.T) {
	e := NewError(ErrEmptyName, "invalid value")
	want := "a name is required: invalid value"
	assert.ErrorIs(t, e, ErrEmptyName)
	assert.Equal(t, e.Error(), want)
}
