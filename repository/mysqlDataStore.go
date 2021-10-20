package repository

import (
	"git-clones/go-api-simple/data"
)

type Repository interface {
	Close()
	GetAllEmployees() ([]data.Employee, error)
	GetEmployeeById(id string) (data.Employee, error)
	CreateEmployee(emp data.Employee) (data.Employee, error)
	DeleteEmployee(id string) error
	UpdateEmployee(emp data.Employee) (data.Employee, error)
}
