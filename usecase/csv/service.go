package csv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vsantosalmeida/csv-parser/entity"
	errs "github.com/vsantosalmeida/csv-parser/pkg/errors"
)

type service struct {
	patterns map[string]*FilePattern
	inMemDB  map[string]string
}

const (
	writeEmployeesFile = "writeEmployeesFile"
	writeBadDataFile   = "writeEmployeesFile"
)

// NewParser returns a Parser interface to process CSV files.
//
// a map[string]*FilePattern is required to translate the columns names for each file.
func NewParser(filePatternMap map[string]*FilePattern) (Parser, error) {
	if err := validateFilePattern(filePatternMap); err != nil {
		return nil, err
	}

	return &service{
		patterns: filePatternMap,
		inMemDB:  make(map[string]string),
	}, nil
}

func (s *service) ParseFiles(files []string) (errors map[string]error) {
	errors = make(map[string]error)
	badDataResult := make(map[string][]*BadData)
	employeesResult := make([]*entity.Employee, 0)

	log.WithFields(log.Fields{
		"event": "processing_files",
		"total": len(files),
		"files": files,
	}).Debug()
	for _, file := range files {
		csvFile, err := os.Open(file)
		if err != nil {
			log.WithFields(log.Fields{
				"event":  "open_file_failed",
				"file":   file,
				"reason": err,
			}).Error("could not open the file")
			errors[file] = errs.NewError(errs.ErrOpeningFile, err.Error())
			continue
		}

		records, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			log.WithFields(log.Fields{
				"event":  "read_file_failed",
				"file":   file,
				"reason": err,
			}).Error("could not read the file in csv format")
			errors[file] = errs.NewError(errs.ErrReadingFile, err.Error())
			continue
		}
		csvFile.Close()

		filePattern, ok := s.patterns[file]
		if !ok {
			log.WithFields(log.Fields{
				"event": "file_pattern_not_found",
				"file":  file,
			}).Error("a file pattern was not found to process the given csv file")
			errors[file] = errs.ErrUnprocessableFile
			continue
		}

		log.WithFields(log.Fields{
			"event": "processing_file",
			"file":  file,
		}).Info()

		employees, badData := s.mapEmployeeOrBadData(records, filePattern)
		employeesResult = append(employeesResult, employees...)
		if len(badData) != 0 {
			log.WithFields(log.Fields{
				"event": "file_processed_with_bad_data",
				"file":  file,
			}).Warn("some lines was not successfully processed")
			badDataResult[file] = badData
		}

		log.WithFields(log.Fields{
			"event": "file_processed",
			"file":  file,
		}).Info("file processed without critical errors")
	}

	if err := s.writeBadDataResultFile(badDataResult); err != nil {
		errors[writeBadDataFile] = errs.NewError(errs.ErrWriteFile, err.Error())
	}

	if err := s.writeEmployeesResultFile(employeesResult); err != nil {
		errors[writeEmployeesFile] = errs.NewError(errs.ErrWriteFile, err.Error())
	}

	log.WithFields(log.Fields{
		"event":               "parse_files_finished",
		"file_errors":         len(errors),
		"employees_processed": len(employeesResult),
	}).Info()

	return
}

func (s *service) mapEmployeeOrBadData(records [][]string, pattern *FilePattern) (employees []*entity.Employee, badData []*BadData) {
	header := make(map[int]string)

	for i, line := range records {
		if i == 0 {
			for j, columnName := range line {
				header[j] = columnName
			}
		} else {
			employeeMap := make(map[string]string)

			for k, value := range line {
				propertyName := header[k]
				employeeMap[propertyName] = value
			}

			log.WithFields(log.Fields{
				"event": "building_new_employee",
				"line":  i + 1,
			}).Info("")
			employee, reasons, ok := s.buildEmployee(employeeMap, pattern)
			if !ok {
				log.WithFields(log.Fields{
					"event": "unprocessable_line",
					"line":  i + 1,
				}).Warn("the line have invalid properties")
				badData = append(badData, &BadData{
					Line:    json.Number(fmt.Sprintf("%d", i+1)),
					Reasons: reasons,
				})
				continue
			}

			employees = append(employees, employee)
		}
	}

	return
}

func (s *service) buildEmployee(employeeMap map[string]string, pattern *FilePattern) (
	employee *entity.Employee, reasons []string, ok bool) {
	name, err := buildAndValidateName(employeeMap[pattern.FirstNameColumn], employeeMap[pattern.LastNameColumn])
	if err != nil {
		log.WithFields(log.Fields{
			"event":  "name_validation_failed",
			"reason": err,
		}).Error("error when validating employee name")
		reasons = append(reasons, err.Error())
	}

	salary, err := buildAndValidateSalary(employeeMap[pattern.SalaryColumn])
	if err != nil {
		log.WithFields(log.Fields{
			"event":  "salary_validation_failed",
			"reason": err,
		}).Error("error when validating employee salary")
		reasons = append(reasons, err.Error())
	}

	email, err := s.trimAndValidateEmail(employeeMap[pattern.EmailColumn])
	if err != nil {
		log.WithFields(log.Fields{
			"event":  "email_validation_failed",
			"reason": err,
		}).Error("error when validating employee e-mail")
		reasons = append(reasons, err.Error())
	}

	id, err := s.validateID(employeeMap[pattern.IDColumn])
	if err != nil {
		log.WithFields(log.Fields{
			"event":  "id_validation_failed",
			"reason": err,
		}).Error("error when validating employee ID")
		reasons = append(reasons, err.Error())
	}

	phone := employeeMap[pattern.PhoneColumn]

	if len(reasons) != 0 {
		return
	}

	employee = &entity.Employee{
		ID:     id,
		Email:  email,
		Name:   name,
		Salary: salary,
		Phone:  phone,
	}

	ok = true

	log.WithFields(log.Fields{
		"event": "employee_created",
		"id":    id,
	}).Info("employee created with success")

	return
}

func (s *service) trimAndValidateEmail(email string) (string, error) {
	email = strings.Trim(email, " ")
	if _, err := mail.ParseAddress(email); err != nil {
		return "", errs.ErrInvalidEmailFormat
	}

	if _, ok := s.inMemDB[email]; ok {
		return "", errs.ErrEmailConstraintViolation
	}

	s.inMemDB[email] = ""

	return email, nil
}

func (s *service) validateID(id string) (string, error) {
	id = strings.Trim(id, " ")
	if id == "" {
		return "", errs.ErrInvalidIDValue
	}

	if _, ok := s.inMemDB[id]; ok {
		return "", errs.ErrIDConstraintViolation
	}

	s.inMemDB[id] = ""

	return id, nil
}

func buildAndValidateName(firstName, lastName string) (string, error) {
	firstName = strings.Trim(firstName, " ")
	lastName = strings.Trim(lastName, " ")
	fNameOk := firstName != ""
	lNameOK := lastName != ""

	if !fNameOk {
		return "", errs.ErrEmptyName
	}

	if fNameOk && lNameOK {
		return entity.BuildEmployeeName(firstName, lastName), nil
	}

	return firstName, nil
}

func buildAndValidateSalary(salary string) (float64, error) {
	salary = strings.Trim(salary, "$,. ")

	floatSalary, err := strconv.ParseFloat(salary, 64)
	if err != nil || floatSalary <= 0 {
		return 0, errs.ErrInvalidSalaryValue
	}

	return floatSalary, nil
}
