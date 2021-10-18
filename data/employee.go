package data

import (
	"encoding/json"
	"strings"
	"time"
)

type CustomDOB time.Time

type Employee struct {
	Id           string    `json:"id"`
	FirstName    string    `json:"firstname" validate:"required"`
	MiddleName   string    `json:"middlename" validate:"required"`
	LastName     string    `json:"lastname" validate:"required"`
	Gender       string    `json:"gender" validate:"required"`
	Salary       float64   `json:"salary" validate:"required"`
	DOB          CustomDOB `json:"dob" validate:"required"`
	Email        string    `json:"email" validate:"required,email"`
	Phone        string    `json:"phone" validate:"required"`
	State        string    `json:"state" validate:"required"`
	Postcode     int       `json:"postcode" validate:"required"`
	AddressLine1 string    `json:"addressline1" validate:"required"`
	AddressLine2 string    `json:"addressline2" validate:"required"`
	TFN          string    `json:"tfn" validate:"required"`
	SuperBalance float64   `json:"superbalance" validate:"required"`
}

// Implement Marshaler and Unmarshaler interface
//Unmarshal JSON -> go
func (j *CustomDOB) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = CustomDOB(t)
	return nil
}

//encode GO values to JSON
func (j CustomDOB) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

// Maybe a Format function for printing your date
func (j CustomDOB) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}
