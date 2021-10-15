package service

import (
	"database/sql"
	"emp-simple/config"
	"emp-simple/data"
	"emp-simple/errorhandling"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	//test := c.Params.ByName("dob")
	c.BindJSON(&emp)

	dob, err := emp.DOB.MarshalJSON()
	log.Println("log:", dob)
	emp.DOB.UnmarshalJSON(dob)
	log.Println("log:", emp.DOB)
	if err != nil {
		c.JSON(500, gin.H{"message": "Marshaling JSON error"})
	}
	t, _ := time.Parse(time.RFC1123Z, string(dob))
	emp.DOB = data.CustomDOB(t)

	_, err = config.Db.Exec("INSERT INTO employee (id,first_name ,middle_name ,last_name ,gender ,salary ,dob ,email , phone , state ,postcode, address_line1 ,address_line2, tfn, super_balance) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", emp.Id, emp.FirstName, emp.MiddleName, emp.LastName, emp.Gender, emp.Salary, emp.DOB, emp.Email, emp.Phone, emp.State, emp.Postcode, emp.AddressLine1, emp.AddressLine2, emp.TFN, emp.SuperBalance)
	if err != nil {
		c.String(http.StatusInternalServerError, "error with inserting")
		return
	}
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	//fmt.Errorf("insert employee: %v", err)
	// }
	c.JSON(http.StatusOK, emp)

}
