package csv

import "github.com/vsantosalmeida/csv-parser/entity"

// Parser is an abstraction to parse CSV files and normalize the data to an *entity.Employee
// for any error on processing a line will be added to a map[string][]*BadData with the file name as the key,
// and a slice with the failed lines.
//
// When the process is finished will write JSON files with successes and failures.
type Parser interface {
	Writer

	// ParseFiles is responsible to read each CSV file from the []string and process it.
	//
	// In case of error to process a file will add the error to map[string]error with the file name as the key
	// and the received error as a value.
	ParseFiles(files []string) (errors map[string]error)
}

// Writer is an embedded interface in Parser, responsible to write files with the results from the
// Parser.ParseFiles method.
type Writer interface {
	writeEmployeesResultFile(employees []*entity.Employee) error
	writeBadDataResultFile(badData map[string][]*BadData) error
}
