package csv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vsantosalmeida/csv-parser/entity"
)

const filenamePrefix = "%s-%s.json"

func (s *service) writeEmployeesResultFile(employees []*entity.Employee) error {
	if len(employees) > 0 {
		file, err := json.MarshalIndent(employees, "", " ")
		if err != nil {
			log.WithFields(log.Fields{
				"event":  "marshal_employees_failed",
				"reason": err,
			}).Error("could not parse Employee to a json structure")
			return err
		}

		err = ioutil.WriteFile(fmt.Sprintf(filenamePrefix, "employee", time.Now().Format("20060102150405")), file, 0644)
		if err != nil {
			log.WithFields(log.Fields{
				"event":  "write_employee_file_failed",
				"reason": err,
			}).Error()
			return err
		}

		log.WithFields(log.Fields{
			"event": "employee_result_file_wrote",
		}).Info()
	}

	return nil
}

func (s *service) writeBadDataResultFile(badData map[string][]*BadData) error {
	if len(badData) > 0 {
		file, err := json.MarshalIndent(badData, "", " ")
		if err != nil {
			log.WithFields(log.Fields{
				"event":  "marshal_bad_data_failed",
				"reason": err,
			}).Error("could not parse BadData to a json structure")
			return err
		}
		err = ioutil.WriteFile(fmt.Sprintf(filenamePrefix, "badData", time.Now().Format("20060102150405")), file, 0644)
		if err != nil {
			log.WithFields(log.Fields{
				"event":  "write_bad_data_file_failed",
				"reason": err,
			}).Error()
			return err
		}

		log.WithFields(log.Fields{
			"event": "bad_data_result_file_wrote",
		}).Info()
	}

	return nil
}
