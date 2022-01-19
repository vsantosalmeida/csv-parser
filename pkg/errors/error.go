package errors

import "fmt"

const (
	ErrEmptyFilePatternMapReceived = err("could not create a parser without a map of FilePattern")
	ErrEmptyName                   = err("a name is required")
	ErrInvalidSalaryValue          = err("could not convert salary to a float value or salary is less or equals to 0")
	ErrInvalidEmailFormat          = err("e-mail must be a valid address ex: email@example.com")
	ErrInvalidIDValue              = err("an id is required")
	ErrInvalidFilePattern          = err("the columns for ID, FirstName, Salary and Email are required to process a file")
	ErrEmailConstraintViolation    = err("e-mail already used by an employee")
	ErrIDConstraintViolation       = err("this ID is already used by an employee")
	ErrOpeningFile                 = err("could not open the given file")
	ErrReadingFile                 = err("could not read the given file")
	ErrUnprocessableFile           = err("could not find a file pattern to process")
	ErrWriteFile                   = err("could not write the result file")
)

type err string

func (e err) Error() string {
	return string(e)
}

// NewError is a syntax sugar for fmt.Errorf("%w: %s", err, cause)
// used to give more context to the err error
func NewError(err err, cause string) error {
	return fmt.Errorf("%w: %s", err, cause)
}
