package mysql

import (
	"database/sql"
	"fmt"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/errorhandling"
	"git-clones/go-api-simple/repository"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type MysqlRepo struct {
	Mysqldb *sql.DB
}

func (r *MysqlRepo) Close() {
	r.Mysqldb.Close()
}

func NewMySQLRepository(dialect string, config mysql.Config, idleConn, maxConn int) (repository.Repository, error) {
	db, err := sql.Open(dialect, "tester:secret@tcp(db:3306)/IGD?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("failed to open sql")
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping error: %s", err)
	}
	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)
	return &MysqlRepo{db}, nil // returing address of MysqlRepo type variable
}

//this will have the implementation of repository interface for CRUD functions for mysql
func executeQuery(q string, r *MysqlRepo) (*sql.Rows, error) {
	rows, err := r.Mysqldb.Query(q)
	if err != nil {
		rows.Close()
		return nil, errorhandling.WrapError("mysql.executeQuery function", errorhandling.BadRequest, "something wrong with query")
	}
	return rows, nil
}

func (r *MysqlRepo) GetAllEmployees() ([]data.Employee, error) {
	var emps []data.Employee
	rows, err := executeQuery("SELECT * FROM employee", r)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var emp data.Employee
		if err := rows.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName,
			&emp.LastName, &emp.Gender, &emp.Salary, &emp.DOB, &emp.Email,
			&emp.Phone, &emp.State, &emp.Postcode, &emp.AddressLine1, &emp.AddressLine2,
			&emp.TFN, &emp.SuperBalance); err != nil {
			return nil, err
		}
		emps = append(emps, emp)
	}
	return emps, nil
}

func (r *MysqlRepo) GetEmployeeById(id string) (data.Employee, error) {
	emp, err := getEmployeeById(id, r)
	if err != nil { //if not found
		return data.Employee{}, err
	}
	return emp, nil //if found
}

func getEmployeeById(id string, r *MysqlRepo) (data.Employee, error) {
	var emp data.Employee
	q := "SELECT * FROM employee WHERE id = ?"
	row := r.Mysqldb.QueryRow(q, id)
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
func (r *MysqlRepo) CreateEmployee(emp data.Employee) (data.Employee, error) {

	emp.Id = uuid.New().String()       // generate a new random UUID and assign
	tm := emp.DOB.Format("2006-01-02") //format into string
	_, err := r.Mysqldb.Exec("INSERT INTO employee (id,first_name ,middle_name ,last_name ,gender ,salary ,dob ,email , phone , state ,postcode, address_line1 ,address_line2, tfn, super_balance) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", emp.Id, emp.FirstName, emp.MiddleName, emp.LastName, emp.Gender, emp.Salary, tm, emp.Email, emp.Phone, emp.State, emp.Postcode, emp.AddressLine1, emp.AddressLine2, emp.TFN, emp.SuperBalance)
	if err != nil {
		return data.Employee{}, err
	}
	return emp, nil
}

func (r *MysqlRepo) DeleteEmployee(id string) error {
	_, err := getEmployeeById(id, r) //check if id exists
	if err != nil {                  //if not exists return
		return fmt.Errorf("sql: no rows in result set")
	}
	query := "DELETE FROM employee WHERE id = ?" //delete emp
	_, err = r.Mysqldb.Exec(query, id)
	if err != nil {
		return &errorhandling.RequestError{Context: "mysql.GetAllEmployeeById delete query execution", Code: errorhandling.Internal, Message: err.Error()}
	}
	return nil //if deletion successful
}

func (r *MysqlRepo) UpdateEmployee(emp data.Employee) (data.Employee, error) {
	var newEmp data.Employee
	//if param is empty, do not update
	q := "UPDATE employee SET first_name = ?,middle_name = ?,last_name = ? ,gender = ?,salary = ?,dob = ?,email = ?, phone = ?, state = ? ,postcode = ?, address_line1 = ?,address_line2 = ?, tfn = ?, super_balance = ? WHERE id = ?"
	originalEmp, err := getEmployeeById(emp.Id, r)
	if err != nil { //if not exists return
		return data.Employee{}, fmt.Errorf("sql: no rows in result set")
	}
	log.Println(originalEmp)
	tm := emp.DOB.Format("2006-01-02") //format into string

	_, err = r.Mysqldb.Exec(q, emp.FirstName, emp.MiddleName,
		emp.LastName, emp.Gender, emp.Salary, tm, emp.Email,
		emp.Phone, emp.State, emp.Postcode, emp.AddressLine1, emp.AddressLine2,
		emp.TFN, emp.SuperBalance, emp.Id)
	if err != nil {
		return data.Employee{}, &errorhandling.RequestError{Context: "mysql.go.UpdateEmployee r.SqlDb.Exec update", Code: errorhandling.Internal, Message: err.Error()}
	}
	newEmp, _ = getEmployeeById(emp.Id, r)
	log.Println("Updated emp: ", newEmp)
	return newEmp, nil
}
