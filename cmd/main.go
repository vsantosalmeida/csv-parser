package main

import (
	"flag"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vsantosalmeida/csv-parser/usecase/csv"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var f string
	flag.StringVar(&f, "f", "", `Files names separated by ","`)
	flag.Parse()

	if strings.Trim(f, " ") == "" {
		log.WithFields(log.Fields{
			"event": "empty_files_arg",
		}).Panic("the `-f` arg is required to process a file and must not be empty")
	}

	files := strings.Split(f, ",")

	filePatterns := csv.NewFilePatternMap(files)

	parser, err := csv.NewParser(filePatterns)
	if err != nil {
		log.WithFields(log.Fields{
			"event":  "create_csv_parser_error",
			"reason": err,
		}).Panic("could not create a parser with given configurations")
	}

	errs := parser.ParseFiles(files)
	if len(errs) != 0 {
		for k, v := range errs {
			log.WithFields(log.Fields{
				"event":     "parse_file_finished_with_errors",
				"error_key": k,
				"reason":    v,
			}).Warn()
		}
		return
	}

	log.WithFields(log.Fields{
		"event": "all_files_processed",
		"files": files,
	}).Info()
}
