package csv

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsantosalmeida/csv-parser/entity"
	"github.com/vsantosalmeida/csv-parser/pkg/errors"
)

func TestService_writeEmployeesResultFile(t *testing.T) {
	var (
		givenEmployees = []*entity.Employee{
			{
				ID:     "1",
				Email:  "doe@test.com",
				Name:   "John Doe",
				Salary: 10,
			},
			{
				ID:     "2",
				Email:  "Mary@tes.com",
				Name:   "Mary Jane",
				Salary: 15,
			},
			{
				ID:     "3",
				Email:  "max@test.com",
				Name:   "Max Topperson",
				Salary: 11,
			},
		}
		employeePattern = "*employee*.json"
	)

	svc := service{}
	err := svc.writeEmployeesResultFile(givenEmployees)
	assert.NoError(t, err)

	eMatches, err := filepath.Glob(employeePattern)
	if err != nil {
		t.FailNow()
	}
	assert.NotEmpty(t, eMatches)
	deleteFile(eMatches[0], t)
}

func TestService_writeBadDataResultFile(t *testing.T) {
	var (
		givenBadData = map[string][]*BadData{
			"file.csv": {
				{
					Line:    "2",
					Reasons: []string{errors.ErrEmptyName.Error(), errors.ErrInvalidSalaryValue.Error()},
				},
				{
					Line:    "3",
					Reasons: []string{errors.ErrInvalidEmailFormat.Error(), errors.ErrInvalidIDValue.Error()},
				},
				{
					Line:    "4",
					Reasons: []string{errors.ErrEmptyName.Error(), errors.ErrInvalidEmailFormat.Error()},
				},
				{
					Line:    "6",
					Reasons: []string{errors.ErrInvalidSalaryValue.Error(), errors.ErrEmailConstraintViolation.Error(), errors.ErrIDConstraintViolation.Error()},
				},
			},
		}
		badDataPattern = "*badData*.json"
	)

	svc := service{}
	err := svc.writeBadDataResultFile(givenBadData)
	assert.NoError(t, err)

	eMatches, err := filepath.Glob(badDataPattern)
	if err != nil {
		t.FailNow()
	}
	assert.NotEmpty(t, eMatches)
	deleteFile(eMatches[0], t)
}

func TestService_writeEmployeesResultFile_JsonError(t *testing.T) {
	var (
		givenEmployees = []*entity.Employee{
			{
				ID:     "1",
				Email:  "doe@test.com",
				Name:   "John Doe",
				Salary: math.Inf(10),
			},
			{
				ID:     "2",
				Email:  "Mary@tes.com",
				Name:   "Mary Jane",
				Salary: 15,
			},
			{
				ID:     "3",
				Email:  "max@test.com",
				Name:   "Max Topperson",
				Salary: 11,
			},
		}
		employeePattern = "*employee*.json"
	)

	svc := service{}
	err := svc.writeEmployeesResultFile(givenEmployees)
	assert.Error(t, err)

	eMatches, err := filepath.Glob(employeePattern)
	if err != nil {
		t.FailNow()
	}
	assert.Empty(t, eMatches)
}

func TestService_writeBadDataResultFile_JsonError(t *testing.T) {
	var (
		givenBadData = map[string][]*BadData{
			"file.csv": {
				{
					Line:    "2",
					Reasons: []string{errors.ErrEmptyName.Error(), errors.ErrInvalidSalaryValue.Error()},
				},
				{
					Line:    "3",
					Reasons: []string{errors.ErrInvalidEmailFormat.Error(), errors.ErrInvalidIDValue.Error()},
				},
				{
					Line:    "4",
					Reasons: []string{errors.ErrEmptyName.Error(), errors.ErrInvalidEmailFormat.Error()},
				},
				{
					Line:    "invalid",
					Reasons: []string{errors.ErrInvalidSalaryValue.Error(), errors.ErrEmailConstraintViolation.Error(), errors.ErrIDConstraintViolation.Error()},
				},
			},
		}
		badDataPattern = "*badData*.json"
	)

	svc := service{}
	err := svc.writeBadDataResultFile(givenBadData)
	assert.Error(t, err)

	eMatches, err := filepath.Glob(badDataPattern)
	if err != nil {
		t.FailNow()
	}
	assert.Empty(t, eMatches)
}

func deleteFile(file string, t *testing.T) {
	err := os.Remove(file)
	if err != nil {
		t.Fatalf("delete file: %s error: %q", file, err)
	}
}
