package data

import (
	"encoding/json"
	"strings"
	"time"
)

type CustomDOB time.Time

type Employee struct {
	Id           string    `json:"id"`
	FirstName    string    `json:"firstname"`
	MiddleName   string    `json:"middlename"`
	LastName     string    `json:"lastname"`
	Gender       string    `json:"gender"`
	Salary       float64   `json:"salary"`
	DOB          CustomDOB `json:"dob"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	State        string    `json:"state"`
	Postcode     int       `json:"postcode"`
	AddressLine1 string    `json:"addressline1"`
	AddressLine2 string    `json:"addressline2"`
	TFN          string    `json:"tfn"`
	SuperBalance float64   `json:"superbalance"`
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
