package repository

import (
	"git-clones/go-api-simple/data"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	GetAllEmployees(c *gin.Context) ([]data.Employee, error)
	GetEmployeeById(c *gin.Context) (data.Employee, error)
	CreateEmployee(c *gin.Context) error
	DeleteEmployee(c *gin.Context) error
	UpdateEmployee(c *gin.Context) error
}
