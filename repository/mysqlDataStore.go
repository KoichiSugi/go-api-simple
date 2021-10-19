package repository

import (
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/errorhandling"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	GetAllEmployees(c *gin.Context) ([]data.Employee, errorhandling.RequestError)
	GetEmployeeById(c *gin.Context) (data.Employee, error)
	CreateEmployee(c *gin.Context) error
	DeleteEmployee(c *gin.Context) error
	UpdateEmployee(c *gin.Context) error
}
