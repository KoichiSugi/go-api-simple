package service_test

import (
	"database/sql"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/service"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var c *gin.Context
var format = "2006-01-02"
var mockDOB, _ = time.Parse(format, "2018-12-12")

var u = &data.Employee{
	Id:           "123",
	FirstName:    "Kyoko",
	MiddleName:   "B",
	LastName:     "Fukada",
	Gender:       "Female",
	Salary:       9999.99,
	DOB:          data.CustomDOB(mockDOB),
	Email:        "momo@gmail.com",
	Phone:        "03999999",
	State:        "VIC",
	Postcode:     1234,
	AddressLine1: "ABC street 123",
	AddressLine2: "JPN",
	TFN:          "123456",
	SuperBalance: 100.0,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestExecuteQueryWithId(t *testing.T) {
	//db, mock := NewMock()
	query := "SELECT * FROM employee"
	rows := sqlmock.NewRows([]string{"Id", "firstname", "middlename", "lastname", "gender", "salary", "dob", "email", "phone", "state", "postcode", "addressline1", "addressline2", "tfn", "superbalance"}).
		AddRow(u.Id, u.FirstName, u.MiddleName, u.LastName, u.Gender, u.Salary, "1994-01-01", u.Email, u.Phone, u.State, u.Postcode, u.AddressLine1, u.AddressLine2, u.TFN, u.SuperBalance)

	mock.ExpectQuery(query).WithArgs(u.Id).WillReturnRows(rows)
	result := service.ExecuteQuery(u.Id, c)
	log.Print(result)
	assert.NotNil(t, result)

}
