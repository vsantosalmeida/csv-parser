package csv

import "encoding/json"

type BadData struct {
	Line    json.Number `json:"line"`
	Reasons []string    `json:"reasons"`
}
