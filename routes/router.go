package routes

import (
	"emp-simple/service"
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
	router := gin.Default()
	group1 := router.Group("/employee")
	{
		group1.GET("/", service.GetAllEmployees)
		group1.GET("/:id", service.GetEmployeeById)
		group1.POST("/", service.CreateEmployee)
	}
	return router
}
