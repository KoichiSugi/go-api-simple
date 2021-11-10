package dto

import (
	"time"
)

//DTO object which is going to be interacting/ passed around on client side
type CustomDOB time.Time
type Employee struct {
	Id           string    `json:"id" db:"id"`
	FirstName    string    `json:"firstname" validate:"required" db:"first_name"`
	MiddleName   string    `json:"middlename" validate:"max=10,excludesall=!()#@{}" db:"middle_name"`
	LastName     string    `json:"lastname" validate:"required" db:"last_name"`
	Gender       string    `json:"gender" validate:"required" db:"gender"`
	Salary       float64   `json:"salary" validate:"required" db:"salary"`
	DOB          CustomDOB `json:"dob" validate:"required" db:"dob"`
	Email        string    `json:"email" validate:"required,email" db:"email"`
	Phone        string    `json:"phone" validate:"required" db:"phone"`
	State        string    `json:"state" validate:"required" db:"state"`
	Postcode     int       `json:"postcode" validate:"required" db:"postcode"`
	AddressLine1 string    `json:"addressline1" validate:"required" db:"address_line1"`
	AddressLine2 string    `json:"addressline2" validate:"required" db:"address_line2"`
	TFN          string    `json:"tfn" validate:"required" db:"tfn"`
	SuperBalance float64   `json:"superbalance" validate:"required" db:"super_balance"`
}
