package mysql

import (
	"database/sql"
	"fmt"
	"git-clones/go-api-simple/config"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/errorhandling"
	"git-clones/go-api-simple/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type MysqlRepo struct {
	SqlDb *sql.DB
}

func (r *MysqlRepo) Close() {
	r.SqlDb.Close()
}

func NewMySQLRepository(dialect string, config mysql.Config, idleConn, maxConn int) (repository.Repository, error) {
	db, err := sql.Open(dialect, config.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)
	return &MysqlRepo{db}, nil
}

//this will have the implementation of repository interface for CRUD functions for mysql
func executeQuery(q string) (*sql.Rows, error) {
	rows, err := config.Db.Query(q)
	if err != nil {
		//c.JSON(500, gin.H{"message": "Something is wrong with query or db"})
		rows.Close()
		return nil, errorhandling.WrapError("mysql.executeQuery function", errorhandling.BadRequest, "something wrong with query")
	}
	return rows, nil
}

func (r *MysqlRepo) GetAllEmployees(c *gin.Context) ([]data.Employee, error) {
	var emps []data.Employee
	rows, err := executeQuery("SELECT * FROM employee")
	if err != nil {
		c.JSON(int(errorhandling.BadRequest), err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var emp data.Employee
		if err := rows.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName,
			&emp.LastName, &emp.Gender, &emp.Salary, &emp.DOB, &emp.Email,
			&emp.Phone, &emp.State, &emp.Postcode, &emp.AddressLine1, &emp.AddressLine2,
			&emp.TFN, &emp.SuperBalance); err != nil {
			c.JSON(500, errorhandling.WrapError("mysql.GetAllEmployee rows.scan", errorhandling.BadRequest, err.Error()))
			return nil, err
		}
		emps = append(emps, emp)
		c.JSON(http.StatusOK, emp)
	}
	return emps, nil
}

func (r *MysqlRepo) GetEmployeeById(c *gin.Context) error {
	id := c.Params.ByName("id")
	emp, err := getEmployeeById(id)
	if err != nil { //if not found
		c.JSON(404, errorhandling.WrapError("mysql.GetEmployeeByIdHandler sql.ErrNoRows", errorhandling.NotFound, err.Error()))
		return err
	}
	//if found
	c.JSON(http.StatusOK, emp)
	return nil
}

func getEmployeeById(id string) (data.Employee, error) {
	var emp data.Employee
	q := "SELECT * FROM employee WHERE id = ?"
	row := config.Db.QueryRow(q, id)
	if err := row.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName,
		&emp.LastName, &emp.Gender, &emp.Salary, &emp.DOB, &emp.Email,
		&emp.Phone, &emp.State, &emp.Postcode, &emp.AddressLine1, &emp.AddressLine2,
		&emp.TFN, &emp.SuperBalance); err != nil {
		if err == sql.ErrNoRows {
			return data.Employee{}, err
		}
		return data.Employee{}, err
	}
	return emp, nil
}
func (r *MysqlRepo) CreateEmployee(c *gin.Context) error {
	var emp data.Employee
	v := validator.New()

	if err := c.BindJSON(&emp); err != nil {
		err := v.Struct(emp)
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
		}
		c.JSON(int(errorhandling.BadRequest), &errorhandling.RequestError{Context: " CreateEmployee c.BindJson", Code: errorhandling.BadRequest, Message: err.Error()})
		return err
	}
	emp.Id = uuid.New().String() // generate a new random UUID and assign
	log.Println(emp)
	tm := emp.DOB.Format("2006-01-02") //format into string
	log.Println(tm)
	_, err := config.Db.Exec("INSERT INTO employee (id,first_name ,middle_name ,last_name ,gender ,salary ,dob ,email , phone , state ,postcode, address_line1 ,address_line2, tfn, super_balance) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", emp.Id, emp.FirstName, emp.MiddleName, emp.LastName, emp.Gender, emp.Salary, tm, emp.Email, emp.Phone, emp.State, emp.Postcode, emp.AddressLine1, emp.AddressLine2, emp.TFN, emp.SuperBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errorhandling.RequestError{Context: "insert confif.Db.Exec", Code: errorhandling.BadRequest, Message: err.Error()})
		return err
	}
	c.JSON(http.StatusOK, emp)
	return nil
}

func (r *MysqlRepo) DeleteEmployee(c *gin.Context) error {
	id := c.Params.ByName("id")
	//find if id exists
	_, err := getEmployeeById(id)
	if err != nil { //if not exists return
		c.JSON(404, errorhandling.WrapError("mysql.GetAllEmployeeById sql.ErrNoRows", errorhandling.NotFound, err.Error()))
		return err
	}
	//delete emp
	query := "DELETE FROM employee WHERE id = ?"
	_, err = config.Db.Exec(query, id)
	if err != nil {
		c.JSON(500, errorhandling.WrapError("mysql.GetAllEmployeeById Db.exec", errorhandling.Internal, err.Error()))
		return err
	}
	c.String(200, fmt.Sprintf("emp id %v has been deleted", id))
	return nil
}

func (r *MysqlRepo) UpdateEmployee(c *gin.Context) error {
	var originalEmp, newEmp data.Employee
	id := c.Params.ByName("id")
	//if param is empty, do not update
	q := "UPDATE employee SET first_name = ?,middle_name = ?,last_name = ? ,gender = ?,salary = ?,dob = ?,email = ?, phone = ?, state = ? ,postcode = ?, address_line1 = ?,address_line2 = ?, tfn = ?, super_balance = ? WHERE id = ?"
	originalEmp, err := getEmployeeById(id) //check if emp exists
	if err != nil {
		c.JSON(404, errorhandling.WrapError("mysql.go.UpdateEmployee", errorhandling.NotFound, err.Error()))
		return err
	}
	log.Println(originalEmp)
	if err = c.BindJSON(&newEmp); err != nil { //check if it binds correctly
		c.JSON(int(errorhandling.BadRequest), &errorhandling.RequestError{Context: "UpdateEmployee c.BindJson", Code: errorhandling.BadRequest, Message: err.Error()})
		return err
	}
	tm := originalEmp.DOB.Format("2006-01-02") //format into string
	log.Println(tm)
	//if originalEmp.FirstName != newEmp.FirstName
	result, err := config.Db.Exec(q, newEmp.FirstName, newEmp.MiddleName,
		newEmp.LastName, newEmp.Gender, newEmp.Salary, tm, newEmp.Email,
		newEmp.Phone, newEmp.State, newEmp.Postcode, newEmp.AddressLine1, newEmp.AddressLine2,
		newEmp.TFN, newEmp.SuperBalance, id)
	if err != nil {
		c.JSON(404, errorhandling.WrapError("mysql.go.UpdateEmployee config.Db.Exec", errorhandling.NotFound, err.Error()))
	}
	c.JSON(200, result)
	return nil
}
