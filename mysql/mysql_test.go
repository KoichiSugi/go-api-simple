package mysql_test

import (
	"database/sql"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/mysql"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

var u2 = data.Employee{
	Id:           "123234",
	FirstName:    "Kyoko",
	MiddleName:   "B",
	LastName:     "Fukada",
	Gender:       "Female",
	Salary:       9999.99,
	DOB:          data.CustomDOB(mockDOB),
	Email:        "momose@gmail.com",
	Phone:        "01298384",
	State:        "VIC",
	Postcode:     1234,
	AddressLine1: "ABC street 123",
	AddressLine2: "JPN",
	TFN:          "19944991",
	SuperBalance: 100.0,
}
var emp = data.Employee{
	Id:           uuid.New().String(),
	FirstName:    "Vinodh",
	MiddleName:   "K",
	LastName:     "Landa",
	Gender:       "Male",
	Salary:       555.55,
	DOB:          data.CustomDOB(mockDOB),
	Email:        "vinod@gmail.com",
	Phone:        "01298384",
	AddressLine1: "Lonsdale",
	AddressLine2: "street",
	State:        "vic",
	Postcode:     3000,
	TFN:          "19944991",
	SuperBalance: 4645,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	//db, mock, err := sqlmock.New()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestGetAllEmployees(t *testing.T) {
	db, mock := NewMock()
	repo := &mysql.MysqlRepo{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT * FROM employee"
	rows := sqlmock.NewRows([]string{"Id", "firstname", "middlename", "lastname", "gender", "salary", "dob", "email", "phone", "state", "postcode", "addressline1", "addressline2", "tfn", "superbalance"}).
		AddRow(u.Id, u.FirstName, u.MiddleName, u.LastName, u.Gender, u.Salary, mockDOB, u.Email, u.Phone, u.State, u.Postcode, u.AddressLine1, u.AddressLine2, u.TFN, u.SuperBalance)

	mock.ExpectQuery(query).WillReturnRows(rows)
	users, err := repo.GetAllEmployees()
	log.Println("users: ", users)
	log.Println("error", err)
	assert.NotEmpty(t, users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestGetEmployeeById(t *testing.T) {
	db, mock := NewMock()
	repo := &mysql.MysqlRepo{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT * FROM employee WHERE id = ?"
	rows := sqlmock.NewRows([]string{"Id", "firstname", "middlename", "lastname", "gender", "salary", "dob", "email", "phone", "state", "postcode", "addressline1", "addressline2", "tfn", "superbalance"}).
		AddRow(u.Id, u.FirstName, u.MiddleName, u.LastName, u.Gender, u.Salary, mockDOB, u.Email, u.Phone, u.State, u.Postcode, u.AddressLine1, u.AddressLine2, u.TFN, u.SuperBalance)

	mock.ExpectQuery(query).WillReturnRows(rows)
	user, err := repo.GetEmployeeById(u.Id)
	log.Println("user returned: ", user)
	log.Println("error", err)
	assert.NotNil(t, user)
	assert.NoError(t, err)
}

func TestGetEmployeeByIdFailure(t *testing.T) {
	db, mock := NewMock()
	repo := &mysql.MysqlRepo{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT * FROM employee WHERE id = ?"
	rows := sqlmock.NewRows([]string{"Id", "firstname", "middlename", "lastname", "gender", "salary", "dob", "email", "phone", "state", "postcode", "addressline1", "addressline2", "tfn", "superbalance"}).
		AddRow(u.Id, u.FirstName, u.MiddleName, u.LastName, u.Gender, u.Salary, mockDOB, u.Email, u.Phone, u.State, u.Postcode, u.AddressLine1, u.AddressLine2, u.TFN, u.SuperBalance)

	mock.ExpectQuery(query).WithArgs(u.Id).WillReturnRows(rows)

	user, err := repo.GetEmployeeById("nonExistingId")
	log.Println("user returned: ", user)
	log.Println("error", err)
	assert.Empty(t, user)
	assert.Error(t, err)
}

func TestCreateEmployee(t *testing.T) { //error
	db, mock := NewMock()
	repo := &mysql.MysqlRepo{db}
	defer repo.Close()
	query := "INSERT INTO employee"
	mock.MatchExpectationsInOrder(false)
	mock.ExpectExec(query).WithArgs(emp.Id, emp.FirstName, emp.MiddleName, emp.LastName, emp.Gender, emp.Salary, emp.DOB, emp.Email, emp.Phone,
		emp.AddressLine1, emp.AddressLine2, emp.State, emp.Postcode, emp.TFN, emp.SuperBalance).WillReturnResult(sqlmock.NewResult(0, 1))

	employee, err := repo.CreateEmployee(emp)
	assert.NotEmpty(t, employee)
	log.Println("empl and err values", employee, err)
	assert.NoError(t, err)

}

func TestCreateEmployeeFailure(t *testing.T) {

}

func TestDeleteEmployee(t *testing.T) { //error
	db, mock := NewMock()
	repo := &mysql.MysqlRepo{db}
	defer repo.Close()
	query := "Delete from employee"
	mock.ExpectExec(query).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
	err := repo.DeleteEmployee(emp.Id)
	log.Println("The deleted row with no error", err, string(emp.Id))
	assert.NoError(t, err)
}
