package service

import (
	"database/sql"
	"fmt"
	"git-clones/go-api-simple/config"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/errorhandling"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

func ExecuteQuery(q string, c *gin.Context) *sql.Rows {
	rows, err := config.Db.Query(q)
	if err != nil {
		c.JSON(500, gin.H{"message": "Something is wrong with query or db"})
		rows.Close()
		log.Fatal("Something wrong with query", err)
	}
	return rows
}

func ExecuteQueryWithId(q string, c *gin.Context, id string) *sql.Row {
	rows := config.Db.QueryRow(q, id)
	return rows
}

func GetAllEmployees(c *gin.Context) {
	context := "service.GetAllEmployees"
	rows := ExecuteQuery("SELECT * FROM employee", c)
	defer rows.Close()
	for rows.Next() {
		var emp data.Employee
		if err := rows.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName,
			&emp.LastName, &emp.Gender, &emp.Salary, &emp.DOB, &emp.Email,
			&emp.Phone, &emp.State, &emp.Postcode, &emp.AddressLine1, &emp.AddressLine2,
			&emp.TFN, &emp.SuperBalance); err != nil {
			log.Fatal(err)
			//test
			c.JSON(500, &errorhandling.RequestError{Context: context, Code: errorhandling.Internal, Message: "DB is empty"})
			//c.JSON(500, gin.H{"message": "DB is empty"})
		}
		c.JSON(http.StatusOK, emp)
	}
}

func GetEmployeeById(c *gin.Context) {
	var emp data.Employee
	id := c.Params.ByName("id")
	query := "SELECT * FROM employee WHERE id = ?"
	row := ExecuteQueryWithId(query, c, id)
	if err := row.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName,
		&emp.LastName, &emp.Gender, &emp.Salary, &emp.DOB, &emp.Email,
		&emp.Phone, &emp.State, &emp.Postcode, &emp.AddressLine1, &emp.AddressLine2,
		&emp.TFN, &emp.SuperBalance); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "employee not found "})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "database error "})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func CreateEmployee(c *gin.Context) {
	var emp data.Employee
	v := validator.New()

	if err := c.BindJSON(&emp); err != nil {
		err := v.Struct(emp)
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
		}
		c.JSON(int(errorhandling.BadRequest), &errorhandling.RequestError{Context: " CreateEmployee c.BindJson", Code: errorhandling.BadRequest, Message: err.Error()})
		return
	}
	emp.Id = uuid.New().String() // generate a new random UUID and assign
	log.Println(emp)
	tm := emp.DOB.Format("2006-01-02") //format into string
	log.Println(tm)
	_, err := config.Db.Exec("INSERT INTO employee (id,first_name ,middle_name ,last_name ,gender ,salary ,dob ,email , phone , state ,postcode, address_line1 ,address_line2, tfn, super_balance) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", emp.Id, emp.FirstName, emp.MiddleName, emp.LastName, emp.Gender, emp.Salary, tm, emp.Email, emp.Phone, emp.State, emp.Postcode, emp.AddressLine1, emp.AddressLine2, emp.TFN, emp.SuperBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errorhandling.RequestError{Context: "insert confif.Db.Exec", Code: errorhandling.BadRequest, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, emp)

}

func DeleteEmployee(c *gin.Context) {
	id := c.Params.ByName("id")
	var emp, err = ExpGetEmployeeById(id)
	if err.Message != "" {

	}
	if (data.Employee{}) == emp {
		fmt.Println("is zero value")
	}
}

func ExpExecuteQueryWithId(q, id string) *sql.Row {
	rows := config.Db.QueryRow(q, id)
	return rows
}

func ExpGetEmployeeById(id string) (data.Employee, errorhandling.RequestError) {
	var emp data.Employee
	query := "SELECT * FROM employee WHERE id = ?"
	row := ExpExecuteQueryWithId(query, id)
	if err := row.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName,
		&emp.LastName, &emp.Gender, &emp.Salary, &emp.DOB, &emp.Email,
		&emp.Phone, &emp.State, &emp.Postcode, &emp.AddressLine1, &emp.AddressLine2,
		&emp.TFN, &emp.SuperBalance); err != nil {
		if err == sql.ErrNoRows {
			return data.Employee{}, errorhandling.RequestError{Context: "sql.ErrNoRows", Code: errorhandling.Internal, Message: err.Error()}
		}
		return data.Employee{}, errorhandling.RequestError{Context: "ExpGetEmployeeById", Code: errorhandling.Internal, Message: err.Error()}
	}
	return emp, errorhandling.RequestError{}
}

// func ValidateEmployee(emp data.Employee) errorhandling.RequestError {

// }
