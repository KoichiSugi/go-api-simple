package mysql

import (
	"fmt"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/errs"
	"git-clones/go-api-simple/repository"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MysqlRepo struct {
	Mysqldb *sqlx.DB
}

func (r *MysqlRepo) Close() {
	r.Mysqldb.Close()
}

func NewMySQLRepository(dialect string, config mysql.Config, idleConn, maxConn int) (repository.Repository, error) {
	db, err := sqlx.Open(dialect, "tester:secret@tcp(db:3306)/IGD?parseTime=true")
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

func (r *MysqlRepo) GetAllEmployees() ([]data.Employee, error) {
	emps := make([]data.Employee, 0)
	findAllSql := "SELECT * FROM employee"
	err := r.Mysqldb.Select(&emps, findAllSql)
	if err != nil {
		return nil, err
	}
	return emps, nil
}

func (r *MysqlRepo) GetEmployeeById(id string) (data.Employee, error) {
	var e data.Employee
	getEmpByIdSql := "SELECT * FROM employee WHERE id = ?"
	err := r.Mysqldb.Get(&e, getEmpByIdSql, id)
	if err != nil { //if not found
		return data.Employee{}, err
	}
	return e, nil //if found
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
	_, err := r.GetEmployeeById(id) //check if id exists
	if err != nil {                 //if not exists return
		return fmt.Errorf("sql: no rows in result set")
	}
	query := "DELETE FROM employee WHERE id = ?" //delete emp
	_, err = r.Mysqldb.Exec(query, id)
	if err != nil {
		return &errs.RequestError{Context: "mysql.GetAllEmployeeById delete query execution", Code: errs.Internal, Message: err.Error()}
	}
	return nil //if deletion successful
}

func (r *MysqlRepo) UpdateEmployee(emp data.Employee) (data.Employee, error) {
	var newEmp data.Employee
	//if param is empty, do not update
	q := "UPDATE employee SET first_name = ?,middle_name = ?,last_name = ? ,gender = ?,salary = ?,dob = ?,email = ?, phone = ?, state = ? ,postcode = ?, address_line1 = ?,address_line2 = ?, tfn = ?, super_balance = ? WHERE id = ?"
	originalEmp, err := r.GetEmployeeById(emp.Id)
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
		return data.Employee{}, &errs.RequestError{Context: "mysql.go.UpdateEmployee r.SqlDb.Exec update", Code: errs.Internal, Message: err.Error()}
	}
	newEmp, _ = r.GetEmployeeById(emp.Id)
	log.Println("Updated emp: ", newEmp)
	return newEmp, nil
}
