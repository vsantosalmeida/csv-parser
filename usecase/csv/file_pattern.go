package csv

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	errs "github.com/vsantosalmeida/csv-parser/pkg/errors"
)

// FilePattern it's a struct to translate the columns from a CSV file to map as an entity.Employee.
type FilePattern struct {
	FirstNameColumn string
	LastNameColumn  string
	SalaryColumn    string
	EmailColumn     string
	IDColumn        string
	PhoneColumn     string
}

// NewFilePatternMap with each file in the []string will build a *FilePattern to process the received CSV file,
// when called will receive as input all the columns' names.
func NewFilePatternMap(files []string) map[string]*FilePattern {
	patterns := make(map[string]*FilePattern)
	for _, file := range files {
		reader := bufio.NewReader(os.Stdin)
		log.WithFields(log.Fields{
			"event": "creating_file_pattern",
			"file":  file,
		}).Info("")

		fmt.Println("If file dont have the column name, just hit enter")
		fmt.Println("Enter First Name column name:")
		firstName, _ := reader.ReadString('\n')
		firstName = strings.TrimSuffix(firstName, "\n")

		fmt.Println("Enter Last Name column name:")
		lastName, _ := reader.ReadString('\n')
		lastName = strings.TrimSuffix(lastName, "\n")

		fmt.Println("Enter Salary column name:")
		salary, _ := reader.ReadString('\n')
		salary = strings.TrimSuffix(salary, "\n")

		fmt.Println("Enter Email column name:")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSuffix(email, "\n")

		fmt.Println("Enter ID column name:")
		id, _ := reader.ReadString('\n')
		id = strings.TrimSuffix(id, "\n")

		fmt.Println("Enter Phone column name:")
		phone, _ := reader.ReadString('\n')
		phone = strings.TrimSuffix(phone, "\n")

		patterns[file] = &FilePattern{
			FirstNameColumn: firstName,
			LastNameColumn:  lastName,
			SalaryColumn:    salary,
			EmailColumn:     email,
			IDColumn:        id,
			PhoneColumn:     phone,
		}
	}

	return patterns
}

func validateFilePattern(filePatternMap map[string]*FilePattern) error {
	if len(filePatternMap) == 0 {
		return errs.ErrEmptyFilePatternMapReceived
	}

	for fileName, filePattern := range filePatternMap {
		if filePattern.EmailColumn == "" || filePattern.IDColumn == "" || filePattern.SalaryColumn == "" ||
			filePattern.FirstNameColumn == "" {
			return errs.NewError(errs.ErrInvalidFilePattern, fileName)
		}
	}

	return nil
}
