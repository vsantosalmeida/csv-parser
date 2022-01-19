package entity

import "fmt"

type Employee struct {
	ID     string  `json:"id"`
	Email  string  `json:"email"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Phone  string  `json:"phone,omitempty"`
}

func BuildEmployeeName(firstName, lastName string) string {
	return fmt.Sprintf("%s %s", firstName, lastName)
}
