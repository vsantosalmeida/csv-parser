package csv_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vsantosalmeida/csv-parser/entity"
	"github.com/vsantosalmeida/csv-parser/pkg/errors"
	"github.com/vsantosalmeida/csv-parser/usecase/csv"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestNewParser(t *testing.T) {
	var givenFilePatterns = map[string]*csv.FilePattern{
		"file.csv": {
			FirstNameColumn: "Name",
			SalaryColumn:    "Wage",
			EmailColumn:     "Email",
			IDColumn:        "Number",
		},
	}

	svc, err := csv.NewParser(givenFilePatterns)
	assert.NotNil(t, svc)
	assert.NoError(t, err)
}

func TestNewParser_Error(t *testing.T) {
	tt := []struct {
		name              string
		givenFilePatterns map[string]*csv.FilePattern
		wantErr           error
	}{
		{
			name: "Without FirstName Column",
			givenFilePatterns: map[string]*csv.FilePattern{
				"file.csv": {
					SalaryColumn: "Wage",
					EmailColumn:  "Email",
					IDColumn:     "Number",
				},
			},
			wantErr: errors.ErrInvalidFilePattern,
		},
		{
			name: "Without Salary Column",
			givenFilePatterns: map[string]*csv.FilePattern{
				"file.csv": {
					FirstNameColumn: "Name",
					EmailColumn:     "Email",
					IDColumn:        "Number",
				},
			},
			wantErr: errors.ErrInvalidFilePattern,
		},
		{
			name: "Without Email Column",
			givenFilePatterns: map[string]*csv.FilePattern{
				"file.csv": {
					FirstNameColumn: "Name",
					SalaryColumn:    "Wage",
					IDColumn:        "Number",
				},
			},
			wantErr: errors.ErrInvalidFilePattern,
		},
		{
			name: "Without ID Column",
			givenFilePatterns: map[string]*csv.FilePattern{
				"file.csv": {
					FirstNameColumn: "Name",
					EmailColumn:     "Email",
					SalaryColumn:    "Wage",
				},
			},
			wantErr: errors.ErrInvalidFilePattern,
		},
		{
			name:    "Empty FilePattern Map",
			wantErr: errors.ErrEmptyFilePatternMapReceived,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := csv.NewParser(tc.givenFilePatterns)
			assert.Nil(t, svc)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestService_ParseFiles_ExampleFile1(t *testing.T) {
	var (
		givenFile = "test_files/roster1.csv"

		givenFilePatterns = map[string]*csv.FilePattern{
			givenFile: {
				FirstNameColumn: "Name",
				SalaryColumn:    "Wage",
				EmailColumn:     "Email",
				IDColumn:        "Number",
			},
		}

		files = []string{givenFile}

		wantEmployees = []*entity.Employee{
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

		wantBadData = map[string][]*csv.BadData{
			givenFile: {
				{
					Line:    "5",
					Reasons: []string{errors.ErrInvalidEmailFormat.Error()},
				},
				{
					Line:    "6",
					Reasons: []string{errors.ErrEmailConstraintViolation.Error()},
				},
			},
		}
	)

	svc, err := csv.NewParser(givenFilePatterns)
	assert.NoError(t, err)

	errs := svc.ParseFiles(files)
	gotEmployees, gotBadData, files := getResults(t)
	assert.Equal(t, wantEmployees, gotEmployees)
	assert.Equal(t, wantBadData, gotBadData)
	assert.Empty(t, errs)
	deleteFiles(files, t)
}

//
func TestService_ParseFiles_ExampleFile2(t *testing.T) {
	var (
		givenFile = "test_files/roster2.csv"

		givenFilePatterns = map[string]*csv.FilePattern{
			givenFile: {
				FirstNameColumn: "First",
				LastNameColumn:  "Last",
				SalaryColumn:    "Salary",
				EmailColumn:     "E-mail",
				IDColumn:        "ID",
			},
		}
		files = []string{givenFile}

		wantEmployees = []*entity.Employee{
			{
				ID:     "RT1",
				Email:  "doe@test.com",
				Name:   "John Doe",
				Salary: 10,
			},
			{
				ID:     "RT2",
				Email:  "mary@tes.com",
				Name:   "Mary Jane",
				Salary: 15,
			},
			{
				ID:     "RT3",
				Email:  "max@test.com",
				Name:   "Max Topperson",
				Salary: 11,
			},
			{
				ID:     "RT4",
				Email:  "alfred@test.com",
				Name:   "Alfred Donald",
				Salary: 11.5,
			},
			{
				ID:     "RT5",
				Email:  "jane.doe@test.com",
				Name:   "Jane Doe",
				Salary: 8.45,
			},
		}
	)

	svc, err := csv.NewParser(givenFilePatterns)
	assert.NoError(t, err)

	errs := svc.ParseFiles(files)
	gotEmployees, gotBadData, files := getResults(t)
	assert.Equal(t, wantEmployees, gotEmployees)
	assert.Empty(t, gotBadData)
	assert.Empty(t, errs)
	deleteFiles(files, t)
}

func TestService_ParseFiles_ExampleFile3(t *testing.T) {
	var (
		givenFile         = "test_files/roster3.csv"
		givenFilePatterns = map[string]*csv.FilePattern{
			givenFile: {
				FirstNameColumn: "first name",
				LastNameColumn:  "last name",
				SalaryColumn:    "Rate",
				EmailColumn:     "e-mail",
				IDColumn:        "Employee Number",
				PhoneColumn:     "Mobile",
			},
		}

		files = []string{givenFile}

		wantEmployees = []*entity.Employee{
			{
				ID:     "RT2",
				Email:  "mary@tes.com",
				Name:   "Mary Jane",
				Salary: 15,
				Phone:  "1448561274",
			},
			{
				ID:     "RT4",
				Email:  "alfred@test.com",
				Name:   "Alfred Donald",
				Salary: 11.5,
				Phone:  "2145385777",
			},
		}

		wantBadData = map[string][]*csv.BadData{
			givenFile: {
				{
					Line:    "2",
					Reasons: []string{errors.ErrInvalidSalaryValue.Error()},
				},
				{
					Line:    "4",
					Reasons: []string{errors.ErrInvalidIDValue.Error()},
				},
				{
					Line:    "6",
					Reasons: []string{errors.ErrInvalidEmailFormat.Error()},
				},
			},
		}
	)

	svc, err := csv.NewParser(givenFilePatterns)
	assert.NoError(t, err)

	errs := svc.ParseFiles(files)
	gotEmployees, gotBadData, files := getResults(t)
	assert.Equal(t, wantEmployees, gotEmployees)
	assert.Equal(t, wantBadData, gotBadData)
	assert.Empty(t, errs)
	deleteFiles(files, t)
}

func TestService_ParseFiles_ExampleFile4(t *testing.T) {
	var (
		givenFile = "test_files/roster4.csv"

		givenFilePatterns = map[string]*csv.FilePattern{
			givenFile: {
				FirstNameColumn: "f. name",
				LastNameColumn:  "l. name",
				SalaryColumn:    "wage",
				EmailColumn:     "email",
				IDColumn:        "emp id",
				PhoneColumn:     "phone",
			},
		}
		files = []string{givenFile}

		wantEmployees = []*entity.Employee{
			{
				ID:     "RT2",
				Email:  "mary@tes.com",
				Name:   "Mary Jane",
				Salary: 15,
				Phone:  "144 856 1274",
			},
			{
				ID:     "RT4",
				Email:  "alfred@test.com",
				Name:   "Alfred Donald",
				Salary: 11.5,
				Phone:  "214 538 5777",
			},
			{
				ID:     "RT5",
				Email:  "jane.doe@test.com",
				Name:   "Jane Doe",
				Salary: 8.45,
			},
		}

		wantBadData = map[string][]*csv.BadData{
			givenFile: {
				{
					Line:    "2",
					Reasons: []string{errors.ErrInvalidSalaryValue.Error()},
				},
				{
					Line:    "4",
					Reasons: []string{errors.ErrInvalidEmailFormat.Error()},
				},
				{
					Line:    "7",
					Reasons: []string{errors.ErrInvalidSalaryValue.Error()},
				},
			},
		}
	)

	svc, err := csv.NewParser(givenFilePatterns)
	assert.NoError(t, err)

	errs := svc.ParseFiles(files)
	gotEmployees, gotBadData, files := getResults(t)
	assert.Equal(t, wantEmployees, gotEmployees)
	assert.Equal(t, wantBadData, gotBadData)
	assert.Empty(t, errs)
	deleteFiles(files, t)
}

func TestService_ParseFiles_ExampleFile5(t *testing.T) {
	var (
		givenFile = "test_files/roster5.csv"

		givenFilePatterns = map[string]*csv.FilePattern{
			givenFile: {
				FirstNameColumn: "f. name",
				LastNameColumn:  "l. name",
				SalaryColumn:    "wage",
				EmailColumn:     "email",
				IDColumn:        "emp id",
				PhoneColumn:     "phone",
			},
		}
		files = []string{givenFile}

		wantEmployees = []*entity.Employee{
			{
				ID:     "RT2",
				Email:  "alfred@test.com",
				Name:   "Alfred Donald",
				Salary: 11,
				Phone:  "214 538 5777",
			},
		}

		wantBadData = map[string][]*csv.BadData{
			givenFile: {
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
	)

	svc, err := csv.NewParser(givenFilePatterns)
	assert.NoError(t, err)

	errs := svc.ParseFiles(files)
	gotEmployees, gotBadData, files := getResults(t)
	assert.Equal(t, wantEmployees, gotEmployees)
	assert.Equal(t, wantBadData, gotBadData)
	assert.Empty(t, errs)
	deleteFiles(files, t)
}

func TestService_ParseFiles_Error(t *testing.T) {
	var (
		givenFilePatterns = map[string]*csv.FilePattern{
			"file.csv": {
				FirstNameColumn: "Name",
				SalaryColumn:    "Wage",
				EmailColumn:     "Email",
				IDColumn:        "Number",
			},
		}
		givenFiles = []string{"not_found.csv", "test_files/bad_file.csv", "test_files/roster1.csv"}
	)

	svc, err := csv.NewParser(givenFilePatterns)
	assert.NotNil(t, svc)
	assert.NoError(t, err)

	errs := svc.ParseFiles(givenFiles)
	assert.ErrorIs(t, errs["not_found.csv"], errors.ErrOpeningFile)
	assert.ErrorIs(t, errs["test_files/bad_file.csv"], errors.ErrReadingFile)
	assert.ErrorIs(t, errs["test_files/roster1.csv"], errors.ErrUnprocessableFile)
}

func getResults(t *testing.T) (employees []*entity.Employee, badData map[string][]*csv.BadData, files []string) {
	t.Helper()
	employeePattern := "*employee*.json"
	badDataPattern := "*badData*.json"

	eMatches, err := filepath.Glob(employeePattern)
	if err != nil {
		t.Fatalf("employee file match error: %q", err)
	}

	bdMatches, err := filepath.Glob(badDataPattern)
	if err != nil {
		t.Fatalf("badData file match error: %q", err)
	}

	if len(eMatches) != 0 {
		eBytes := loadFile(eMatches[0], t)
		err = json.Unmarshal(eBytes, &employees)
		if err != nil {
			t.Fatalf("employee marshal error: %q", err)
		}
		files = append(files, eMatches[0])
	}

	if len(bdMatches) != 0 {
		bdBytes := loadFile(bdMatches[0], t)
		err = json.Unmarshal(bdBytes, &badData)
		if err != nil {
			t.Fatalf("badData marshal error: %q", err)
		}
		files = append(files, bdMatches[0])
	}

	return
}

func deleteFiles(files []string, t *testing.T) {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			t.Fatalf("delete file: %s error: %q", file, err)
		}
	}
}

func loadFile(fileName string, t *testing.T) []byte {
	t.Helper()
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatalf("failed to load file: %s error: %q", fileName, err)
	}
	return b
}
