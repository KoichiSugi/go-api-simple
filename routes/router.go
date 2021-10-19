package routes

import (
	"git-clones/go-api-simple/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type rootHandler func(http.ResponseWriter, *http.Request) error

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r) //call root handler
	if err != nil {
		return
	}
}

func CreateRouter() *gin.Engine {
	var msr mysql.MysqlRepo
	router := gin.Default()
	group1 := router.Group("/employee")
	{
		group1.GET("/", func(ctx *gin.Context) {
			msr.GetAllEmployees(ctx)
		})
		group1.GET("/:id", func(c *gin.Context) {
			msr.GetEmployeeByIdHandler(c)
		})
		group1.POST("/", func(c *gin.Context) {
			msr.CreateEmployee(c)
		})
		group1.DELETE("/:id", func(c *gin.Context) {
			msr.DeleteEmployee(c)
		})
		group1.PUT("/:id", func(c *gin.Context) {
			msr.UpdateEmployee(c)
		})
		//group1.GET("/:id", service.GetEmployeeById)
		//group1.POST("/", service.CreateEmployee)
		//group1.DELETE("/:id", service.DeleteEmployee)

	}
	return router
}
